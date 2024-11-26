package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepositoryDynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepositoryDynamoDB(client *dynamodb.Client) repository.UserRepository {
	return UserRepositoryDynamoDB{
		client:    client,
		tableName: "Users",
	}
}

func (r UserRepositoryDynamoDB) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to marshal user: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r UserRepositoryDynamoDB) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	if len(result.Items) == 0 {
		return entities.User{}, ErrUserNotFound
	}

	var user entities.User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		return entities.User{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return user, nil
}
