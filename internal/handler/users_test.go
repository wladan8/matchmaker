package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.com/matchmaker/internal/matchmaker"

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

func TestMatchmakingHandler_AddUserToPool_OK(t *testing.T) {
	m := matchmaker.New(&matchmaker.Config{
		GroupSize:       5,
		DiffSkill:       100,
		DiffLatency:     100,
		TickerFrequency: 50 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.GatherGroupsProcessing(ctx)
	h := New(m)
	users := generateUsers(1000)
	for _, u := range users {
		user := u
		t.Run("", func(t *testing.T) {
			t.Parallel()
			buf, err := json.Marshal(user)
			assert.NoError(t, err)
			reader := bytes.NewReader(buf)
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users/", reader)
			assert.NotEmpty(t, req)
			h.UserToPool(rw, req)
			assert.Equal(t, http.StatusOK, rw.Code)
		})
	}

	for group := range m.ReceiveGroups() {
		group.CalculateMetrics()
		group.Print()
	}
}
