package repository

import (
	"context"
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

type UserRepositoryDynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepositoryDynamoDB(client *dynamodb.Client) repository.UserRepository {
	return &UserRepositoryDynamoDB{
		client:    client,
		tableName: "Users",
	}
}

func (r *UserRepositoryDynamoDB) CreateUser(ctx context.Context, user *entities.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepositoryDynamoDB) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
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
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}
