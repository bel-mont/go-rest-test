package handlers

import "net/http"

// GetLeaderboard handler
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Leaderboard endpoint"))
}
