package main

import (
	"github.com/joho/godotenv"
	"go-rest-test/handlers"
	"log"
	"net/http"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Initialize server mux and register routes
	mux := handlers.SetupRoutes()

	// Start server
	log.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
