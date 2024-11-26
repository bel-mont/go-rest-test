package http

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/adapters/api"
	"go-rest-test/internal/adapters/web"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
)

func InitializeRoutes(router *gin.Engine, playerRepo repository.PlayerRepository, userRepo repository.UserRepository, replayRepo repository.Repository[entities.Replay]) {
	// Add this health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Initialize API handlers
	playerHandler := api.NewPlayerHandler(playerRepo)
	userHandler := api.NewUserHandler(userRepo)

	// Initialize Web handlers
	playerWebHandler := web.NewPlayerWebHandler(playerRepo)
	userWebHandler := web.NewUserWebHandler()
	homeWebHandler := web.NewHomeWebHandler()
	replayWebHandler := web.NewReplayWebHandler(replayRepo)

	// Register routes
	setupAPIRoutes(router, playerHandler, userHandler)
	setupWebRoutes(router, homeWebHandler, userWebHandler, playerWebHandler, replayWebHandler)
}

func setupAPIRoutes(router *gin.Engine, playerHandler api.PlayerHandler, userHandler api.UserHandler) {
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
}

func setupWebRoutes(router *gin.Engine, homeWebHandler web.HomeWebHandler, userWebHandler web.UserWebHandler, playerWebHandler web.PlayerWebHandler, replayWebHandler web.ReplayWebHandler) {
	router.GET("/", homeWebHandler.RenderHome)
	router.GET("/signup", userWebHandler.RenderSignupForm)
	router.GET("/login", userWebHandler.RenderLoginForm)
	router.GET("/players", playerWebHandler.RenderPlayersList)
	router.GET("/replay", replayWebHandler.RenderIndex)
}
