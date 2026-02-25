package main

import (
	"giftFlow/config"
	"giftFlow/handlers"
	"giftFlow/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to Database
	config.ConnectDB()

	// Seed Mock Data if empty
	seedMockData()

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Mock-User")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API Routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/me", handlers.GetMe)
		api.GET("/gifts", handlers.GetGifts)
		api.POST("/applications", handlers.CreateApplication)
		api.GET("/applications", handlers.GetApplications)
		api.PUT("/applications/:id/status", handlers.UpdateApplicationStatus)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func seedMockData() {
	// Seed Users (Optional if we use middleware map, but good for completeness if we switch to DB)
	// Seed Gifts
	handlers.SeedData()
}
