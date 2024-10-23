package handlers

import "net/http"

// SubmitMatch handler
func SubmitMatch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Submit match endpoint"))
}

// Matchmaking handler
func Matchmaking(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Matchmaking endpoint"))
}
