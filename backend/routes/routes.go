package routes

import (
	"giftflow/handlers"
	"giftflow/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Serve static files
	os.MkdirAll("./uploads", 0755)
	r.Static("/uploads", "./uploads")

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
		api.POST("/gifts", handlers.CreateGift)
		api.PUT("/gifts/:id", handlers.UpdateGift)
		api.DELETE("/gifts/:id", handlers.DeleteGift)
		api.POST("/upload", handlers.UploadGiftImage)
		api.POST("/applications", handlers.CreateApplication)
		api.GET("/applications", handlers.GetApplications)
		api.PUT("/applications/:id/status", handlers.UpdateApplicationStatus)
	}

	return r
}
