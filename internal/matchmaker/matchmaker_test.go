package matchmaker

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.com/matchmaker/internal/entity"
)

func randomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

func generateUsers(n int) []entity.User {
	users := make([]entity.User, n)
	for i := 0; i < n; i++ {
		latency := rand.Float64() * 100
		skill := rand.Float64() * 100
		queueTime := randomTimestamp()
		users[i] = entity.User{
			Name:      fmt.Sprintf("%d-latency:%f-skill:%f", i, latency, skill),
			Latency:   latency,
			Skill:     skill,
			QueueTime: queueTime,
		}
	}
	return users
}

func generateEqualUsers(n int) []entity.User {
	users := make([]entity.User, n)
	latency := rand.Float64() * 100
	skill := rand.Float64() * 100
	queueTime := randomTimestamp()
	for i := 0; i < n; i++ {
		users[i] = entity.User{
			Name:      fmt.Sprintf("%d-latency:%f-skill:%f", i, latency, skill),
			Latency:   latency,
			Skill:     skill,
			QueueTime: queueTime,
		}
	}
	return users
}

func TestMatchmaker_ReceiveGroups_allUsersAreEqual(t *testing.T) {
	m := New(nil)
	m.cfg.TickerFrequency = 50 * time.Millisecond
	m.users = generateEqualUsers(100)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	m.GatherGroupsProcessing(ctx)
	numberOfGroups := 0
	for group := range m.ReceiveGroups() {
		group.CalculateMetrics()
		group.Print()
		numberOfGroups++
	}
	assert.Equal(t, 20, numberOfGroups)
}

func TestMatchmaker_ReceiveGroups_firstUserHasVeryBigSkill(t *testing.T) {
	m := New(nil)
	m.cfg.TickerFrequency = 50 * time.Millisecond
	m.users = generateEqualUsers(100)
	m.users[0].Skill = 1_000
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	m.GatherGroupsProcessing(ctx)
	numberOfGroups := 0
	for group := range m.ReceiveGroups() {
		group.CalculateMetrics()
		group.Print()
		numberOfGroups++
	}
	assert.Equal(t, 5, len(m.users))
	assert.Equal(t, 19, numberOfGroups)
}

func TestMatchmaker_ReceiveGroups_lastUserHasVeryBigLatency(t *testing.T) {
	m := New(nil)
	m.cfg.TickerFrequency = 50 * time.Millisecond
	m.users = generateEqualUsers(100)
	m.users[len(m.users)-1].Latency = 1_000_000
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	m.GatherGroupsProcessing(ctx)
	numberOfGroups := 0
	for group := range m.ReceiveGroups() {
		group.CalculateMetrics()
		group.Print()
		numberOfGroups++
	}
	assert.Equal(t, 5, len(m.users))
	assert.Equal(t, 19, numberOfGroups)
}

func TestMatchmaker_ReceiveGroups_lastUserWaitsLongTime(t *testing.T) {
	m := New(nil)
	m.cfg.TickerFrequency = 50 * time.Millisecond
	m.users = generateEqualUsers(6)
	m.users[len(m.users)-1].QueueTime = m.users[len(m.users)-1].QueueTime.Add(-10 * time.Hour)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	m.GatherGroupsProcessing(ctx)
	numberOfGroups := 0
	for group := range m.ReceiveGroups() {
		group.CalculateMetrics()
		group.Print()
		numberOfGroups++
	}
	assert.Equal(t, 1, len(m.users))
	assert.Equal(t, 1, numberOfGroups)
}
