package main

import (
	"go-rest-test/handlers"
	"log"
	"net/http"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Error loading .env file")
	//}

	// InitDb()

	handlers.Setup()

	log.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
