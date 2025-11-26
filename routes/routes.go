package routes

import (
	"ocr-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	uploadController := controllers.NewUploadController()
	transactionController := controllers.NewTransactionController()

	v1 := router.Group("/api/v1")
	{
		// Upload slip (supports multiple files)
		v1.POST("/upload", uploadController.UploadSlip)

		// Transaction CRUD operations
		v1.POST("/transactions", transactionController.Create)
		v1.GET("/transactions", transactionController.GetAll)
		v1.GET("/transactions/:id", transactionController.GetByID)
		v1.PUT("/transactions/:id", transactionController.Update)
		v1.PATCH("/transactions/:id", transactionController.Update)
		v1.DELETE("/transactions/:id", transactionController.Delete)
	}

	// Health check endpoint - handle both GET and HEAD requests
	healthHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "OCR API is running",
		})
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)
}
