package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
)

type UserDynamoRepository struct {
	BaseDynamoRepository[entities.User]
}

func NewUserDynamoRepository(client *dynamodb.Client) repository.UserRepository {
	baseRepo := NewBaseDynamoRepository[entities.User](client, "Users")
	return UserDynamoRepository{baseRepo}
}

func (r UserDynamoRepository) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	users, err := r.QueryByIndex(ctx, "email-index", "email", email)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	if len(users) == 0 {
		return entities.User{}, repository.ErrItemNotFound
	}

	return users[0], nil
}
