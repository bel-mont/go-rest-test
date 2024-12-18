package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/adapters/repository"
	"go-rest-test/internal/core/entities"
	repository2 "go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/database"
	"go-rest-test/internal/infrastructure/http"
	"go-rest-test/internal/infrastructure/storage"
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

	// Initialize S3 client
	s3Client := storage.InitAWSClient()

	// In local environment, create tables and seed data
	if os.Getenv("ENV") == "local" {
		log.Println("Creating tables and seeding data...")
		if err := database.CreateTables(ctx, dynamoClient); err != nil {
			log.Printf("Error creating tables: %v", err)
		}
		if err := database.SeedData(ctx, dynamoClient); err != nil {
			log.Printf("Error seeding data: %v", err)
		}

		// Add buckets
		storage.CreateBuckets(s3Client)
	}

	// Initialize repositories
	//playerRepo, userRepo := initRepositories(dbPool)
	playerRepo, userRepo, replayRepo, multipartUploadRepo := initRepositories(dynamoClient)

	// Initialize router
	router := gin.Default()
	http.InitializeMiddlewares(router)

	// Set up routes
	http.InitializeRoutes(router, playerRepo, userRepo, replayRepo, multipartUploadRepo, s3Client)

	// Serve static files
	router.Static("/web/static", "./web/static")

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

func initRepositories(client *dynamodb.Client) (
	repository2.Repository[entities.Player],
	repository2.UserRepository,
	repository2.Repository[entities.Replay],
	repository2.Repository[entities.MultipartUpload]) {
	//userRepo := repository.NewUserRepositoryPg(dbPool)
	playerRepo := repository.NewPlayerDynamoRepository(client)
	userRepo := repository.NewUserDynamoRepository(client)
	replayRepo := repository.NewReplayDynamoRepository(client)
	multipartUploadRepo := repository.NewMultipartUploadDynamoRepository(client)
	return playerRepo, userRepo, replayRepo, multipartUploadRepo
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
