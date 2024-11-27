package repository

import (
	"context"
	"go-rest-test/internal/core/entities"
)

type UserRepository interface {
	Repository[entities.User]
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
}
