package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes and returns an http.ServeMux with the configured routes.
func SetupRoutes(router *gin.Engine) {
	router.POST("/auth/login", Login)
	router.GET("/leaderboard", GetLeaderboard)
	router.POST("/matches", SubmitMatch)
	router.GET("/player/:id", GetPlayerStats) // Example for dynamic route
	router.POST("/matchmaking", Matchmaking)
	//http.Handle("/", middleware.AuthMiddleware(http.DefaultServeMux))
}
