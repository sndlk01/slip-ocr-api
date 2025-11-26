package controllers

import (
	"net/http"
	"ocr-api/models"
	"ocr-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	service *services.BudgetService
}

func NewBudgetController() *BudgetController {
	return &BudgetController{
		service: services.NewBudgetService(),
	}
}

type CreateBudgetRequest struct {
	Category     string  `json:"category" binding:"required"`
	MonthlyLimit float64 `json:"monthly_limit" binding:"required"`
	Month        int     `json:"month" binding:"required,min=1,max=12"`
	Year         int     `json:"year" binding:"required"`
}

func (c *BudgetController) Create(ctx *gin.Context) {
	var req CreateBudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	budget := &models.Budget{
		Category:     req.Category,
		MonthlyLimit: req.MonthlyLimit,
		Month:        req.Month,
		Year:         req.Year,
	}

	if err := c.service.Create(budget); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Budget created successfully", "budget": budget})
}

func (c *BudgetController) GetAll(ctx *gin.Context) {
	budgets, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

func (c *BudgetController) GetBudgetStatus(ctx *gin.Context) {
	year, _ := strconv.Atoi(ctx.Query("year"))
	month, _ := strconv.Atoi(ctx.Query("month"))

	if year == 0 || month == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "year and month required"})
		return
	}

	statuses, err := c.service.GetBudgetStatus(month, year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"budget_status": statuses})
}

func (c *BudgetController) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
