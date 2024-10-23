package handlers

import (
	http "net/http"
)

func Setup() {
	http.HandleFunc("/auth/login", Login)
	http.HandleFunc("/leaderboard", GetLeaderboard)
	http.HandleFunc("/matches", SubmitMatch)
	http.HandleFunc("/player/", GetPlayerStats)
	http.HandleFunc("/matchmaking", Matchmaking)

	//http.Handle("/", middleware.AuthMiddleware(http.DefaultServeMux))
}
