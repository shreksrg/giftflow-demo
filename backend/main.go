package main

import (
	"giftflow/config"
	"giftflow/initializers"
	"giftflow/routes"
	"log"
	"os"
)

func main() {
	// Connect to Database
	config.ConnectDB()

	// Seed Mock Data if empty
	initializers.SeedMockData()

	// Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
