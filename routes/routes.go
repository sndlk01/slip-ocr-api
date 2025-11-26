package routes

import (
	"ocr-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	uploadController := controllers.NewUploadController()
	transactionController := controllers.NewTransactionController()
	budgetController := controllers.NewBudgetController()
	subscriptionController := controllers.NewSubscriptionController()
	dashboardController := controllers.NewDashboardController()

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

		// Budget management
		v1.POST("/budgets", budgetController.Create)
		v1.GET("/budgets", budgetController.GetAll)
		v1.GET("/budgets/status", budgetController.GetBudgetStatus)
		v1.DELETE("/budgets/:id", budgetController.Delete)

		// Subscription management
		v1.POST("/subscriptions", subscriptionController.Create)
		v1.GET("/subscriptions", subscriptionController.GetAll)
		v1.DELETE("/subscriptions/:id", subscriptionController.Delete)

		// Dashboard & Analytics
		v1.GET("/dashboard/monthly", dashboardController.GetMonthlyTrend)
		v1.GET("/dashboard/yearly", dashboardController.GetYearlyComparison)
		v1.GET("/dashboard/categories", dashboardController.GetCategoryBreakdown)
		v1.GET("/summary/monthly", transactionController.GetMonthlySummary)
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
