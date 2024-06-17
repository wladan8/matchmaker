package entity

import (
	"strings"
	"time"
)

type User struct {
	Name      string    `json:"name"`
	Skill     float64   `json:"skill"`
	Latency   float64   `json:"latency"`
	QueueTime time.Time `json:"-"`
}

type Users []User

func (u Users) String() string {
	users := make([]string, len(u))
	for i, user := range u {
		users[i] = user.Name
	}
	return strings.Join(users, ",")
}
