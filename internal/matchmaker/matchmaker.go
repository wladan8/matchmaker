package matchmaker

import (
	"context"
	"log/slog"
	"math"
	"sort"
	"sync"
	"time"

	"gitlab.com/matchmaker/internal/entity"
)

const (
	defaultGroupSize       = 5
	defaultDiffSkill       = 10
	defaultDiffLatency     = 50
	defaultTickerFrequency = 1 * time.Second
)

type Matchmaker struct {
	users          []entity.User // Pool of users. They are already sorted by time. Because we add user to end of the array.
	cfg            Config
	lastGroupIndex uint64
	groupChan      chan *entity.Group
	mutex          sync.Mutex
}

type Config struct {
	GroupSize       int           // number of users in one each group
	DiffSkill       float64       // difference between skills in each group
	DiffLatency     float64       // difference between latency in each group
	TickerFrequency time.Duration // frequency of gathering group
}

func New(cfg *Config) *Matchmaker {
	if cfg == nil {
		cfg = &Config{
			GroupSize:       defaultGroupSize,
			DiffSkill:       defaultDiffSkill,
			DiffLatency:     defaultDiffLatency,
			TickerFrequency: defaultTickerFrequency,
		}
	}
	return &Matchmaker{
		users:     make([]entity.User, 0, cfg.GroupSize),
		mutex:     sync.Mutex{},
		groupChan: make(chan *entity.Group),
		cfg:       *cfg,
	}
}

func (m *Matchmaker) Put(ctx context.Context, user entity.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	user.QueueTime = time.Now()
	m.users = append(m.users, user)
	return nil
}

func (m *Matchmaker) ReceiveGroups() <-chan *entity.Group {
	return m.groupChan
}

func (m *Matchmaker) GatherGroupsProcessing(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(m.cfg.TickerFrequency)
		for {
			select {
			case <-ctx.Done():
				close(m.groupChan)
				slog.Info("matchmaker.Processing: finished")
				return
			case <-ticker.C:
				m.mutex.Lock()
				if len(m.users) >= m.cfg.GroupSize {
					sort.Slice(m.users, func(i, j int) bool {
						return m.users[i].QueueTime.Before(m.users[j].QueueTime)
					})
					for j := 0; j < len(m.users); j++ {
						usersInGroup := make([]entity.User, 0, m.cfg.GroupSize)
						usersInGroup = append(usersInGroup, m.users[j])
						m.users = append(m.users[:j], m.users[j+1:]...)

						for i := j; i < len(m.users) && len(usersInGroup) < m.cfg.GroupSize; {
							if math.Abs(usersInGroup[0].Skill-m.users[i].Skill) < m.cfg.DiffSkill &&
								math.Abs(usersInGroup[0].Latency-m.users[i].Latency) < m.cfg.DiffLatency {
								usersInGroup = append(usersInGroup, m.users[i])
								m.users = append(m.users[:i], m.users[i+1:]...)
								continue
							}
							i++
						}
						if len(usersInGroup) == m.cfg.GroupSize { // full group
							j = len(m.users)
							m.lastGroupIndex++
							m.groupChan <- &entity.Group{
								SequenceNumber: m.lastGroupIndex,
								Users:          usersInGroup,
							}
						} else {
							m.users = append(usersInGroup, m.users...)
						}
					}
				}
				m.mutex.Unlock()
			}
		}
	}()
}
