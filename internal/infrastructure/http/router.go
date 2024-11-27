package http

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/adapters/api"
	"go-rest-test/internal/adapters/web"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/middlewares"
)

func InitializeRoutes(router *gin.Engine,
	playerRepo repository.Repository[entities.Player],
	userRepo repository.UserRepository,
	replayRepo repository.Repository[entities.Replay],
	uploadStateRepo repository.Repository[entities.MultipartUpload],
	s3Client *s3.Client) {
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
	replayUploadHandler := api.NewReplayUploadHandler(s3Client, replayRepo, uploadStateRepo)
	replayHandler := api.NewReplayHandler(s3Client, replayRepo)

	// Initialize Web handlers
	playerWebHandler := web.NewPlayerWebHandler(playerRepo)
	userWebHandler := web.NewUserWebHandler()
	homeWebHandler := web.NewHomeWebHandler()
	replayWebHandler := web.NewReplayWebHandler(replayRepo)

	// Register routes
	setupAPIRoutes(router, playerHandler, userHandler, replayUploadHandler, replayHandler)
	setupWebRoutes(router, homeWebHandler, userWebHandler, playerWebHandler, replayWebHandler)
}

func setupAPIRoutes(router *gin.Engine, playerHandler api.PlayerHandler, userHandler api.UserHandler, replayUploadHandler api.ReplayUploadHandler, replayHandler api.ReplayHandler) {
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
		apiGroup.GET("/replay/:id/stream", replayHandler.StreamReplay)
	}
	apiRestrictedGroup := router.Group("/api/restricted")
	apiRestrictedGroup.Use(middlewares.AuthMiddleware())
	{
		apiRestrictedGroup.POST("/replay/upload", replayUploadHandler.UploadHandler)

		// Multipart upload endpoints
		apiRestrictedGroup.POST("/replay/upload/init", replayUploadHandler.InitUploadHandler)
		apiRestrictedGroup.GET("/replay/upload/part-url", replayUploadHandler.GetUploadPartURL)
		apiRestrictedGroup.POST("/replay/upload/complete-part", replayUploadHandler.CompletePart)
		apiRestrictedGroup.POST("/replay/upload/complete", replayUploadHandler.CompleteUpload)

	}
}

func setupWebRoutes(router *gin.Engine, homeWebHandler web.HomeWebHandler, userWebHandler web.UserWebHandler, playerWebHandler web.PlayerWebHandler, replayWebHandler web.ReplayWebHandler) {
	router.GET("/", homeWebHandler.RenderHome)
	router.GET("/signup", userWebHandler.RenderSignupForm)
	router.GET("/login", userWebHandler.RenderLoginForm)
	router.GET("/players", playerWebHandler.RenderPlayersList)
	router.GET("/replay", replayWebHandler.RenderIndex)
	router.GET("/replays/:id", replayWebHandler.RenderViewPage)

	userPagesGroup := router.Group("/u/")
	userPagesGroup.Use(middlewares.AuthMiddleware())
	{
		userPagesGroup.GET("/replay/upload", replayWebHandler.RenderUploadPage)
	}
}
