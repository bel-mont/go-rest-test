package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	repositoryInterface "go-rest-test/internal/core/repository"
)

// BaseDynamoRepository provides generic CRUD operations for DynamoDB
type BaseDynamoRepository[T repositoryInterface.Entity] struct {
	client    *dynamodb.Client
	tableName string
}

// NewBaseDynamoRepository creates a new base repository
func NewBaseDynamoRepository[T repositoryInterface.Entity](client *dynamodb.Client, tableName string) BaseDynamoRepository[T] {
	return BaseDynamoRepository[T]{
		client:    client,
		tableName: tableName,
	}
}

// Create creates a new item
func (r BaseDynamoRepository[T]) Create(ctx context.Context, item T) (T, error) {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return item, fmt.Errorf("failed to marshal item: %w", err)
	}

	// Generate and set UUID directly in the attribute map
	if av["id"].(*types.AttributeValueMemberS).Value == "" {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return item, fmt.Errorf("failed to generate UUID: %w", err)
		}
		av["id"] = &types.AttributeValueMemberS{Value: newUUID.String()}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return item, fmt.Errorf("failed to create item: %w", err)
	}

	// Unmarshal the updated item with the new ID
	var newItem T
	if err := attributevalue.UnmarshalMap(av, &newItem); err != nil {
		return item, fmt.Errorf("failed to unmarshal updated item: %w", err)
	}

	return newItem, nil
}

// Get retrieves an item by ID
func (r BaseDynamoRepository[T]) Get(ctx context.Context, id string) (T, error) {
	var empty T
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return empty, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return empty, repositoryInterface.ErrItemNotFound
	}

	var item T
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return empty, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return item, nil
}

// Update updates an existing item
func (r BaseDynamoRepository[T]) Update(ctx context.Context, item T) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// Delete removes an item by ID
func (r BaseDynamoRepository[T]) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// QueryByIndex performs a query using a GSI
func (r BaseDynamoRepository[T]) QueryByIndex(ctx context.Context, indexName, keyName, keyValue string) ([]T, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :%s", keyName, keyName)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			fmt.Sprintf(":%s", keyName): &types.AttributeValueMemberS{Value: keyValue},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	var items []T
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return items, nil
}

// Scan retrieves all items from the table
func (r BaseDynamoRepository[T]) Scan(ctx context.Context) ([]T, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items: %w", err)
	}

	var items []T
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return items, nil
}
