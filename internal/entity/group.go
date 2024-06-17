package entity

import (
	"fmt"
	"log/slog"
	"time"
)

type Group struct {
	SequenceNumber uint64         `json:"sequence_number"`
	SkillGroup     MinMaxAvgFloat `json:"skill_group"`   // skill in group
	LatencyGroup   MinMaxAvgFloat `json:"latency_group"` // latency in group
	TimeSpend      MinMaxAvgTime  `json:"time_spend"`    // time spent in queue
	Users          Users          `json:"users"`         // list of users
}

func (g *Group) CalculateMetrics() {
	if len(g.Users) > 0 {
		maxLatency, minLatency := g.Users[0].Latency, g.Users[0].Latency
		maxSkill, minSkill := g.Users[0].Skill, g.Users[0].Skill
		maxQueueTime, minQueueTime := time.Since(g.Users[0].QueueTime), time.Since(g.Users[0].QueueTime)
		var sumLatency, sumSkill float64
		var sumDurationQueue time.Duration
		for _, u := range g.Users {
			if u.Latency > maxLatency {
				maxLatency = u.Latency
			}
			if u.Latency < minLatency {
				minLatency = u.Latency
			}
			if u.Skill > maxSkill {
				maxSkill = u.Skill
			}
			if u.Skill < minSkill {
				minSkill = u.Skill
			}

			if time.Since(u.QueueTime) > maxQueueTime {
				maxQueueTime = time.Since(u.QueueTime)
			}
			if time.Since(u.QueueTime) < minQueueTime {
				minQueueTime = time.Since(u.QueueTime)
			}

			sumLatency += u.Latency
			sumSkill += u.Skill
			sumDurationQueue += time.Since(u.QueueTime)
		}
		numberOfUsers := float64(len(g.Users))
		g.LatencyGroup = MinMaxAvgFloat{
			max: maxLatency,
			min: minLatency,
			avg: sumLatency / numberOfUsers,
		}
		g.SkillGroup = MinMaxAvgFloat{
			max: maxSkill,
			min: minSkill,
			avg: sumSkill / numberOfUsers,
		}
		g.TimeSpend = MinMaxAvgTime{
			max: maxQueueTime,
			min: minQueueTime,
			avg: sumDurationQueue / time.Duration(len(g.Users)),
		}
	}
}

func (g *Group) Print() {
	slog.Info(fmt.Sprintf("sequence_number:\t%d\n"+
		"skill_group:\t%s\n"+
		"latency_group:\t%s\n"+
		"time_spend:\t%s\n"+
		"user_names:\t%s\n",
		g.SequenceNumber, g.SkillGroup, g.LatencyGroup, g.TimeSpend, g.Users))
}

type MinMaxAvgFloat struct {
	min float64
	max float64
	avg float64
}

func (m MinMaxAvgFloat) String() string {
	return fmt.Sprintf("min:%f\tmax:%f\tavg:%f", m.min, m.max, m.avg)
}

type MinMaxAvgTime struct {
	min time.Duration
	max time.Duration
	avg time.Duration
}

func (m MinMaxAvgTime) String() string {
	return fmt.Sprintf("min:%s\tmax:%s\tavg:%s", m.min, m.max, m.avg)
}
