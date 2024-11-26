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

	//"github.com/aws/aws-lambda-go/events"
	//ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

//var ginLambda *ginadapter.GinLambda

func run() {
	// Load environment variables
	utils.LoadEnv()
	log.Printf("Starting run in environment '%s'", os.Getenv("ENV"))

	// Initialize database connection
	//dbPool := database.InitDB()
	//defer dbPool.Close()

	// Initialize DynamoDB client
	ctx := context.Background()
	dynamoClient := database.InitDynamoDB(ctx)

	// In local environment, create tables and seed data
	if os.Getenv("ENV") == "local" {
		log.Println("Creating tables and seeding data...")
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
	//ginLambda = ginadapter.New(router)

	// Start server
	if err := router.Run(":80"); err != nil {
		log.Printf("Failed to start server on port 80. Try running with elevated privileges (sudo/administrator)")
		// Fallback to 8080 if 80 fails
		log.Printf("Falling back to port 8080...")
		log.Fatal(router.Run(":8080"))
	}
}

func initRepositories(client *dynamodb.Client) (repository2.PlayerRepository, repository2.UserRepository) { //playerRepo := repository.NewPlayerRepositoryPg(dbPool)
	//userRepo := repository.NewUserRepositoryPg(dbPool)
	playerRepo := repository.NewPlayerRepositoryDynamoDB(client)
	userRepo := repository.NewUserRepositoryDynamoDB(client)
	return playerRepo, userRepo
}

//func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	log.Printf("Received request: Method=%s, Path=%s", req.HTTPMethod, req.Path)
//	log.Printf("Request headers: %+v", req.Headers)
//
//	resp, err := ginLambda.ProxyWithContext(ctx, req)
//	if err != nil {
//		log.Printf("Error handling request: %v", err)
//	}
//	log.Printf("Response: StatusCode=%d", resp.StatusCode)
//
//	return resp, err
//}

func main() {
	//lambda.Start(Handler)
	run()
}
