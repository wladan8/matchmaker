package handler

import (
	"context"

	"gitlab.com/matchmaker/internal/entity"
)

type MatchmakerInterface interface {
	Put(ctx context.Context, u entity.User) error
	ReceiveGroups() <-chan *entity.Group
}
