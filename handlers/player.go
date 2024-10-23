package handlers

import "net/http"

// GetPlayerStats handler
func GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Player stats endpoint"))
}
