package repository

import (
	"context"
	"go-rest-test/internal/core/entities"
)

type ReplayRepository interface {
	CreateReplay(ctx context.Context, replay entities.Replay) (entities.Replay, error)
	GetReplay(ctx context.Context, id string) (entities.Replay, error)
	GetReplaysByUserID(ctx context.Context, userID string) ([]entities.Replay, error)
	UpdateReplay(ctx context.Context, replay entities.Replay) error
	DeleteReplay(ctx context.Context, id string) error
}
