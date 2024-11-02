package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-test/handlers/api"
	"go-rest-test/handlers/html"
)

// SetupRoutes initializes and returns an http.ServeMux with the configured routes.
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	// Load templates from the "views" folder
	router.LoadHTMLGlob("views/**/*")

	// HTML Pages
	router.GET("/", html.Home)
	router.GET("/auth/login", html.LoginPage)
	router.GET("/matches", html.MatchListPage)
	router.GET("/matches/favorites", html.FavoritesPage)
	router.GET("/matches/submit", html.SubmitMatchPage)
	router.GET("/matches/:id", html.MatchDetailPage)
	router.GET("/leaderboard", html.LeaderboardPage)
	html.SetDB(db)

	// API Endpoints (JSON)
	//router.POST("/auth/login", html.Login)
	//router.GET("/api/players", html.GetPlayers)
	//router.GET("/api/leaderboard", html.GetLeaderboard)
	//router.GET("/api/matches", html.GetMatches)
	//router.POST("/api/matches", html.SubmitMatch)
	//router.GET("/api/player/:id", html.GetPlayerStats)
	//router.POST("/api/matchmaking", html.Matchmaking)
	api.SetDB(db)
}
