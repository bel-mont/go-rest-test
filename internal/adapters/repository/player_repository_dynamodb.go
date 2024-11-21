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

// PlayerRepositoryDynamoDB handles player data operations using DynamoDB.
// It leverages an injected *dynamodb.Client for executing DynamoDB requests.
// The tableName defines the specific DynamoDB table used for player data management.
type PlayerRepositoryDynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewPlayerRepositoryDynamoDB(client *dynamodb.Client) repository.PlayerRepository {
	return &PlayerRepositoryDynamoDB{
		client:    client,
		tableName: "Players",
	}
}

// GetPlayerByID retrieves a player from the DynamoDB table by ID.
func (r *PlayerRepositoryDynamoDB) GetPlayerByID(ctx context.Context, id string) (entities.Player, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return entities.Player{}, fmt.Errorf("dynamodb error: %w", err)
	}

	if result.Item == nil {
		return entities.Player{}, fmt.Errorf("player not found")
	}

	var player entities.Player
	if err := attributevalue.UnmarshalMap(result.Item, &player); err != nil {
		return entities.Player{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return player, nil
}

// UpdatePlayer updates an existing player record in the DynamoDB table using the provided player entity.
func (r *PlayerRepositoryDynamoDB) UpdatePlayer(ctx context.Context, player entities.Player) error {
	av, err := attributevalue.MarshalMap(player)
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	return nil
}

// DeletePlayer deletes a player from DynamoDB by player ID. Returns an error if the deletion fails.
func (r *PlayerRepositoryDynamoDB) DeletePlayer(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	return nil
}

// GetAllPlayers retrieves all players from the DynamoDB table.
func (r *PlayerRepositoryDynamoDB) GetAllPlayers(ctx context.Context) ([]entities.Player, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan players: %w", err)
	}

	var players []entities.Player
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	return players, nil
}

// CreatePlayer creates a new player record in DynamoDB and returns the player ID or an error.
func (r *PlayerRepositoryDynamoDB) CreatePlayer(ctx context.Context, player entities.Player) (string, error) {
	player.ID = uuid.New().String()
	player.LastLogin = time.Now()

	av, err := attributevalue.MarshalMap(player)
	if err != nil {
		return "", fmt.Errorf("failed to marshal player: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create player: %w", err)
	}

	return player.ID, nil
}
