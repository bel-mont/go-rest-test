package repository

import (
	"context"
	"go-rest-test/internal/core/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entities.User) (entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
}
