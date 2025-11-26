package services

import (
	"fmt"
	"ocr-api/config"
	"ocr-api/models"
)

type BudgetService struct{}

func NewBudgetService() *BudgetService {
	return &BudgetService{}
}

func (s *BudgetService) Create(budget *models.Budget) error {
	// Check if budget already exists for this category/month/year
	var existing models.Budget
	result := config.DB.Where("category = ? AND month = ? AND year = ?",
		budget.Category, budget.Month, budget.Year).First(&existing)

	if result.Error == nil {
		return fmt.Errorf("budget already exists for %s in %d/%d", budget.Category, budget.Month, budget.Year)
	}

	result = config.DB.Create(budget)
	if result.Error != nil {
		return fmt.Errorf("failed to create budget: %w", result.Error)
	}
	return nil
}

func (s *BudgetService) GetAll() ([]models.Budget, error) {
	var budgets []models.Budget
	result := config.DB.Order("year DESC, month DESC, category ASC").Find(&budgets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", result.Error)
	}
	return budgets, nil
}

func (s *BudgetService) GetByMonthYear(month, year int) ([]models.Budget, error) {
	var budgets []models.Budget
	result := config.DB.Where("month = ? AND year = ?", month, year).
		Order("category ASC").Find(&budgets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", result.Error)
	}
	return budgets, nil
}

func (s *BudgetService) GetByID(id uint) (*models.Budget, error) {
	var budget models.Budget
	result := config.DB.First(&budget, id)
	if result.Error != nil {
		return nil, fmt.Errorf("budget not found: %w", result.Error)
	}
	return &budget, nil
}

func (s *BudgetService) Update(id uint, updates map[string]interface{}) (*models.Budget, error) {
	var budget models.Budget
	result := config.DB.First(&budget, id)
	if result.Error != nil {
		return nil, fmt.Errorf("budget not found: %w", result.Error)
	}

	result = config.DB.Model(&budget).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update budget: %w", result.Error)
	}

	return &budget, nil
}

func (s *BudgetService) Delete(id uint) error {
	result := config.DB.Delete(&models.Budget{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("budget not found")
	}
	return nil
}

type BudgetStatus struct {
	Category     string  `json:"category"`
	MonthlyLimit float64 `json:"monthly_limit"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"`
	PercentUsed  float64 `json:"percent_used"`
	Status       string  `json:"status"` // ok, warning, exceeded
}

func (s *BudgetService) GetBudgetStatus(month, year int) ([]BudgetStatus, error) {
	budgets, err := s.GetByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var statuses []BudgetStatus

	for _, budget := range budgets {
		// Calculate spent for this category in this month
		var spent float64
		monthStr := fmt.Sprintf("%02d", month)
		yearStr := fmt.Sprintf("%d", year)
		pattern := fmt.Sprintf("%%/%s/%s", monthStr, yearStr)

		config.DB.Model(&models.Transaction{}).
			Where("type = ? AND category = ? AND date LIKE ?", "expense", budget.Category, pattern).
			Select("COALESCE(SUM(amount), 0)").Scan(&spent)

		remaining := budget.MonthlyLimit - spent
		percentUsed := 0.0
		if budget.MonthlyLimit > 0 {
			percentUsed = (spent / budget.MonthlyLimit) * 100
		}

		status := "ok"
		if percentUsed >= 100 {
			status = "exceeded"
		} else if percentUsed >= 80 {
			status = "warning"
		}

		statuses = append(statuses, BudgetStatus{
			Category:     budget.Category,
			MonthlyLimit: budget.MonthlyLimit,
			Spent:        spent,
			Remaining:    remaining,
			PercentUsed:  percentUsed,
			Status:       status,
		})
	}

	return statuses, nil
}
