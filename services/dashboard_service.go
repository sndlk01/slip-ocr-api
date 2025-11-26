package services

import (
	"fmt"
	"ocr-api/config"
	"ocr-api/models"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

type MonthlyData struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type YearlyData struct {
	Year    int     `json:"year"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type CategoryData struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

// GetMonthlyTrend returns income/expense for 12 months
func (s *DashboardService) GetMonthlyTrend(year int) ([]MonthlyData, error) {
	var data []MonthlyData

	for month := 1; month <= 12; month++ {
		monthStr := fmt.Sprintf("%02d", month)
		yearStr := fmt.Sprintf("%d", year)
		pattern := fmt.Sprintf("%%/%s/%s", monthStr, yearStr)

		var income, expense float64
		config.DB.Model(&models.Transaction{}).
			Where("type = ? AND date LIKE ?", "income", pattern).
			Select("COALESCE(SUM(amount), 0)").Scan(&income)

		config.DB.Model(&models.Transaction{}).
			Where("type = ? AND date LIKE ?", "expense", pattern).
			Select("COALESCE(SUM(amount), 0)").Scan(&expense)

		data = append(data, MonthlyData{
			Month:   fmt.Sprintf("%s/%s", monthStr, yearStr),
			Income:  income,
			Expense: expense,
		})
	}

	return data, nil
}

// GetYearlyComparison compares multiple years
func (s *DashboardService) GetYearlyComparison(years []int) ([]YearlyData, error) {
	var data []YearlyData

	for _, year := range years {
		yearStr := fmt.Sprintf("%d", year)
		pattern := fmt.Sprintf("%%/%s", yearStr)

		var income, expense float64
		config.DB.Model(&models.Transaction{}).
			Where("type = ? AND date LIKE ?", "income", pattern).
			Select("COALESCE(SUM(amount), 0)").Scan(&income)

		config.DB.Model(&models.Transaction{}).
			Where("type = ? AND date LIKE ?", "expense", pattern).
			Select("COALESCE(SUM(amount), 0)").Scan(&expense)

		data = append(data, YearlyData{
			Year:    year,
			Income:  income,
			Expense: expense,
		})
	}

	return data, nil
}

// GetCategoryBreakdown returns spending by category (for pie chart)
func (s *DashboardService) GetCategoryBreakdown(year int, month int, transactionType string) ([]CategoryData, error) {
	monthStr := fmt.Sprintf("%02d", month)
	yearStr := fmt.Sprintf("%d", year)
	pattern := fmt.Sprintf("%%/%s/%s", monthStr, yearStr)

	var data []CategoryData
	result := config.DB.Model(&models.Transaction{}).
		Select("category, SUM(amount) as amount, COUNT(*) as count").
		Where("type = ? AND date LIKE ?", transactionType, pattern).
		Group("category").
		Order("amount DESC").
		Scan(&data)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get category breakdown: %w", result.Error)
	}

	// Replace empty category
	for i := range data {
		if data[i].Category == "" {
			data[i].Category = "ไม่มีหมวดหมู่"
		}
	}

	return data, nil
}
