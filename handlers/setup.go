package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// SetupRoutes initializes and returns an http.ServeMux with the configured routes.
func SetupRoutes(router *gin.Engine) {
	// Load templates from the "views" folder
	router.LoadHTMLGlob("views/*")

	router.GET("/", Home)
	router.GET("/players", GetPlayers) // New route to retrieve all players
	router.POST("/auth/login", Login)
	router.GET("/leaderboard", GetLeaderboard)
	router.GET("/matches", GetMatches)
	router.POST("/matches", SubmitMatch)
	router.GET("/player/:id", GetPlayerStats) // Example for dynamic route
	router.POST("/matchmaking", Matchmaking)
}

// SetDatabase sets the database connection, initialized in main.go
func SetDatabase(database *pgxpool.Pool) {
	db = database
}
