package database

import (
	"context"
	"go-rest-test/pkg/env"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// InitDynamoDB initializes the appropriate DynamoDB client based on the environment
func InitDynamoDB(ctx context.Context) *dynamodb.Client {
	if os.Getenv("ENV") == env.EnvLocal {
		return InitLocalDynamoDB(ctx)
	}
	return InitProductionDynamoDB(ctx)
}

// InitLocalDynamoDB initializes a DynamoDB client for local development
func InitLocalDynamoDB(ctx context.Context) *dynamodb.Client {
	// Basic configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("local"),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	if err != nil {
		log.Fatalf("Failed to load local config: %v", err)
	}

	// Create DynamoDB client with local endpoint
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT"))
	})

	// Test connection
	if _, err := client.ListTables(ctx, &dynamodb.ListTablesInput{}); err != nil {
		log.Fatalf("Failed to connect to local DynamoDB: %v", err)
	}

	log.Println("Successfully initialized local DynamoDB client")
	return client
}

// InitProductionDynamoDB initializes a DynamoDB client for production
func InitProductionDynamoDB(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	// Test connection
	if _, err := client.ListTables(ctx, &dynamodb.ListTablesInput{}); err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	}

	log.Println("Successfully initialized production DynamoDB client")
	return client
}

// CreateTables creates all required DynamoDB tables
func CreateTables(ctx context.Context, client *dynamodb.Client) error {
	tables := []struct {
		name       string
		attributes []types.AttributeDefinition
		keySchema  []types.KeySchemaElement
		gsi        []types.GlobalSecondaryIndex
	}{
		{
			name: "Players",
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("username"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			gsi: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("username-index"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("username"),
							KeyType:       types.KeyTypeHash,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
		},
		{
			name: "Users",
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("email"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			gsi: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("email-index"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("email"),
							KeyType:       types.KeyTypeHash,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
		},
	}

	for _, table := range tables {
		// Check if table exists
		exists, err := tableExists(ctx, client, table.name)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("Table %s already exists", table.name)
			continue
		}

		input := &dynamodb.CreateTableInput{
			AttributeDefinitions:   table.attributes,
			KeySchema:              table.keySchema,
			TableName:              aws.String(table.name),
			GlobalSecondaryIndexes: table.gsi,
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		}

		if _, err := client.CreateTable(ctx, input); err != nil {
			return err
		}
		log.Printf("Created table %s", table.name)
	}
	return nil
}

// tableExists checks if a table exists in DynamoDB
func tableExists(ctx context.Context, client *dynamodb.Client, tableName string) (bool, error) {
	tables, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return false, err
	}

	for _, name := range tables.TableNames {
		if name == tableName {
			return true, nil
		}
	}
	return false, nil
}

// SeedData seeds initial data into the tables
func SeedData(ctx context.Context, client *dynamodb.Client) error {
	// Sample data for Players
	players := []map[string]types.AttributeValue{
		{
			"id":            &types.AttributeValueMemberS{Value: "1"},
			"username":      &types.AttributeValueMemberS{Value: "player1"},
			"level":         &types.AttributeValueMemberN{Value: "1"},
			"total_matches": &types.AttributeValueMemberN{Value: "5"},
			"total_wins":    &types.AttributeValueMemberN{Value: "3"},
		},
		{
			"id":            &types.AttributeValueMemberS{Value: "2"},
			"username":      &types.AttributeValueMemberS{Value: "player2"},
			"level":         &types.AttributeValueMemberN{Value: "2"},
			"total_matches": &types.AttributeValueMemberN{Value: "10"},
			"total_wins":    &types.AttributeValueMemberN{Value: "5"},
		},
	}

	// Only seed data in local environment
	if os.Getenv("ENV") != "local" {
		log.Println("Skipping data seeding in non-local environment")
		return nil
	}

	for _, player := range players {
		input := &dynamodb.PutItemInput{
			TableName: aws.String("Players"),
			Item:      player,
		}

		if _, err := client.PutItem(ctx, input); err != nil {
			return err
		}
	}

	log.Println("Successfully seeded initial data")
	return nil
}
