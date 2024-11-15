package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-test/internal/adapters/repository"
	repository2 "go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/database"
	"go-rest-test/internal/infrastructure/http"
	utils "go-rest-test/pkg/util"
	"log"
	"os"
)

func main() {
	// Load environment variables
	utils.LoadEnv()

	// Initialize database connection
	dbPool := database.InitDB()
	defer dbPool.Close()

	// Initialize repositories
	playerRepo, userRepo := initRepositories(dbPool)

	// Initialize router
	router := gin.Default()
	http.InitializeMiddlewares(router)

	// Set up routes
	http.InitializeRoutes(router, playerRepo, userRepo)

	// Start the server
	startServer(router)
}

func initRepositories(dbPool *pgxpool.Pool) (repository2.PlayerRepository, repository2.UserRepository) {
	playerRepo := repository.NewPlayerRepositoryPg(dbPool)
	userRepo := repository.NewUserRepositoryPg(dbPool)
	return playerRepo, userRepo
}

func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started at port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
