package main

import (
	"log"
	"ocr-api/config"
	"ocr-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	config.Init()

	// Initialize database
	config.InitDatabase()

	// Set Gin mode (release/debug)
	gin.SetMode(gin.DebugMode)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	serverAddr := ":" + config.AppConfig.ServerPort
	log.Printf("Starting server on %s", serverAddr)
	log.Printf("Upload directory: %s", config.AppConfig.UploadDir)
	log.Printf("Database path: %s", config.AppConfig.DatabasePath)
	log.Printf("Tesseract language: %s", config.AppConfig.TesseractLang)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
