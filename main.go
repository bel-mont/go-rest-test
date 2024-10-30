package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-rest-test/db"
	"go-rest-test/handlers"
	"log"
	"os"
	"strings"
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
	setTrustedProxies(router)

	// Setup routes using Gin router
	handlers.SetupRoutes(router)
	// Pass the database connection to handlers
	handlers.SetDatabase(pgxObj)

	// Start server
	log.Println("Server starting at http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}

func setTrustedProxies(router *gin.Engine) {
	// Set trusted proxies from environment variable
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		if err := router.SetTrustedProxies(proxies); err != nil {
			log.Fatalf("Error setting trusted proxies: %v", err)
		}
	} else {
		log.Println("No trusted proxies set; all proxies are trusted by default.")
	}
}
