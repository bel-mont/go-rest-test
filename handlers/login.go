package handlers

import "net/http"

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login endpoint"))
}
