package handlers

import (
	http "net/http"
)

// SetupRoutes initializes and returns an http.ServeMux with the configured routes.
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register each endpoint with the appropriate handler
	mux.HandleFunc("/auth/login", Login)
	mux.HandleFunc("/leaderboard", GetLeaderboard)
	mux.HandleFunc("/matches", SubmitMatch)
	mux.HandleFunc("/player/", GetPlayerStats)
	mux.HandleFunc("/matchmaking", Matchmaking)

	//http.Handle("/", middleware.AuthMiddleware(http.DefaultServeMux))
	return mux
}
