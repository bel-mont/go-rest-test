package repository

import (
	"context"
	"go-rest-test/internal/core/entities"
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player entities.Player) (string, error)
	GetPlayerByID(ctx context.Context, id string) (entities.Player, error)
	UpdatePlayer(ctx context.Context, player entities.Player) error
	DeletePlayer(ctx context.Context, id string) error
	GetAllPlayers(ctx context.Context) ([]entities.Player, error)
}
