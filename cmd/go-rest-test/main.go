package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-rest-test/internal/adapters/api"
	"go-rest-test/internal/adapters/repository"
	"go-rest-test/internal/adapters/web"
	"go-rest-test/internal/infrastructure/database"
	"log"
	"os"
	"strings"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	// Initialize the database connection
	pgxPool := database.InitDB()
	defer pgxPool.Close() // Ensure database connection is closed when the app exits

	// Initialize repository
	// Set up repositories
	playerRepo := repository.NewPlayerRepositoryPg(pgxPool)
	userRepo := repository.NewUserRepositoryPg(pgxPool)

	// Initialize Gin router
	router := gin.Default()
	setTrustedProxies(router)

	// Set up handlers
	playerHandler := api.NewPlayerHandler(playerRepo)
	userHandler := api.NewUserHandler(userRepo)
	playerWebHandler := web.NewPlayerWebHandler(playerRepo)
	userWebHandler := web.NewUserWebHandler()
	homeWebHandler := web.NewHomeWebHandler()
	replayWebHandler := web.NewReplayWebHandler()

	// Player routes (API)
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/players", playerHandler.CreatePlayer)
		apiGroup.GET("/players/:id", playerHandler.GetPlayerByID)
		apiGroup.PUT("/players/:id", playerHandler.UpdatePlayer)
		apiGroup.DELETE("/players/:id", playerHandler.DeletePlayer)
		apiGroup.GET("/players", playerHandler.GetAllPlayers)
		apiGroup.POST("/signup", userHandler.Signup)
		apiGroup.POST("/login", userHandler.Login)
		apiGroup.POST("/logout", userHandler.Logout)
	}

	// Player routes (HTML)
	router.GET("/", homeWebHandler.RenderHome)
	router.GET("/signup", userWebHandler.RenderSignupForm)
	router.GET("/login", userWebHandler.RenderLoginForm)
	router.GET("/players", playerWebHandler.RenderPlayersList)
	router.GET("/replay", replayWebHandler.RenderIndex)

	// Set up some basic routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Placeholder for user registration
	router.POST("/register", func(c *gin.Context) {
		// Registration logic will be added here later
		c.JSON(200, gin.H{
			"message": "User registered",
		})
	})

	// Placeholder for match upload
	router.POST("/upload-match", func(c *gin.Context) {
		// Video upload logic will be added here later
		c.JSON(200, gin.H{
			"message": "Match uploaded",
		})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	log.Println("Server started at port " + port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func setTrustedProxies(router *gin.Engine) {
	// Set trusted proxies from environment variable
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		if err := router.SetTrustedProxies(proxies); err != nil {
			log.Fatalf("Error setting trusted proxies: %v", err)
		}
	} else {
		log.Println("No trusted proxies set; all proxies are trusted by default.")
	}
}
