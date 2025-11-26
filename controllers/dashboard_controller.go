package controllers

import (
	"net/http"
	"ocr-api/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service *services.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{service: services.NewDashboardService()}
}

func (c *DashboardController) GetMonthlyTrend(ctx *gin.Context) {
	year, _ := strconv.Atoi(ctx.Query("year"))
	if year == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "year required"})
		return
	}

	data, err := c.service.GetMonthlyTrend(year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"monthly_trend": data})
}

func (c *DashboardController) GetYearlyComparison(ctx *gin.Context) {
	yearsStr := ctx.Query("years") // e.g., "2023,2024,2025"
	if yearsStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "years required (comma-separated)"})
		return
	}

	var years []int
	for _, y := range strings.Split(yearsStr, ",") {
		if year, err := strconv.Atoi(strings.TrimSpace(y)); err == nil {
			years = append(years, year)
		}
	}

	data, err := c.service.GetYearlyComparison(years)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"yearly_comparison": data})
}

func (c *DashboardController) GetCategoryBreakdown(ctx *gin.Context) {
	year, _ := strconv.Atoi(ctx.Query("year"))
	month, _ := strconv.Atoi(ctx.Query("month"))
	transactionType := ctx.Query("type") // income or expense

	if year == 0 || month == 0 || transactionType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "year, month, and type required"})
		return
	}

	data, err := c.service.GetCategoryBreakdown(year, month, transactionType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"category_breakdown": data})
}
