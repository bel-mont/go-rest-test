package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-rest-test/db"
	"go-rest-test/handlers"
	"log"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	// Initialize the database connection
	pgxObj := db.InitDB()

	// Initialize Gin router
	router := gin.Default()

	// Setup routes using Gin router
	handlers.SetupRoutes(router)
	// Pass the database connection to handlers
	handlers.SetDatabase(pgxObj)

	// Start server
	log.Println("Server starting at http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
