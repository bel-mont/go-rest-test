package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/adapters/repository"
	repository2 "go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/database"
	"go-rest-test/internal/infrastructure/http"
	utils "go-rest-test/pkg/util"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Load environment variables
	utils.LoadEnv()

	// Initialize database connection
	//dbPool := database.InitDB()
	//defer dbPool.Close()

	// Initialize DynamoDB client
	ctx := context.Background()
	dynamoClient := database.InitDynamoDB(ctx)

	// In local environment, create tables and seed data
	if os.Getenv("ENV") == "local" {
		if err := database.CreateTables(ctx, dynamoClient); err != nil {
			log.Printf("Error creating tables: %v", err)
		}
		if err := database.SeedData(ctx, dynamoClient); err != nil {
			log.Printf("Error seeding data: %v", err)
		}
	}

	// Initialize repositories
	//playerRepo, userRepo := initRepositories(dbPool)
	playerRepo, userRepo := initRepositories(dynamoClient)

	// Initialize router
	router := gin.Default()
	http.InitializeMiddlewares(router)

	// Set up routes
	http.InitializeRoutes(router, playerRepo, userRepo)

	// Create the adapter
	ginLambda = ginadapter.New(router)
}

func initRepositories(client *dynamodb.Client) (repository2.PlayerRepository, repository2.UserRepository) { //playerRepo := repository.NewPlayerRepositoryPg(dbPool)
	//userRepo := repository.NewUserRepositoryPg(dbPool)
	playerRepo := repository.NewPlayerRepositoryDynamoDB(client)
	userRepo := repository.NewUserRepositoryDynamoDB(client)
	return playerRepo, userRepo
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
