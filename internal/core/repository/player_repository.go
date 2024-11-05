package repository

import (
	"context"
	"go-rest-test/internal/core/entities"
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player entities.Player) (int, error)
	GetPlayerByID(ctx context.Context, id int) (entities.Player, error)
	UpdatePlayer(ctx context.Context, player entities.Player) error
	DeletePlayer(ctx context.Context, id int) error
	GetAllPlayers(ctx context.Context) ([]entities.Player, error)
}
