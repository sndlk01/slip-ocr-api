package services

import (
	"fmt"
	"ocr-api/config"
	"ocr-api/models"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (s *TransactionService) Create(transaction *models.Transaction) error {
	result := config.DB.Create(transaction)
	if result.Error != nil {
		return fmt.Errorf("failed to create transaction: %w", result.Error)
	}
	return nil
}

// CheckDuplicate checks if a similar transaction already exists
func (s *TransactionService) CheckDuplicate(transaction *models.Transaction) (*models.Transaction, error) {
	var existing models.Transaction

	// Check by reference number first (most reliable)
	if transaction.Reference != "" {
		result := config.DB.Where("reference = ? AND reference != ''", transaction.Reference).
			First(&existing)
		if result.Error == nil {
			return &existing, nil
		}
	}

	// Check by amount + date + time + bank (within same minute)
	if transaction.Date != "" && transaction.Time != "" {
		result := config.DB.Where("amount = ? AND date = ? AND time = ? AND bank = ?",
			transaction.Amount, transaction.Date, transaction.Time, transaction.Bank).
			First(&existing)
		if result.Error == nil {
			return &existing, nil
		}
	}

	return nil, nil
}

func (s *TransactionService) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := config.DB.Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}
	return transactions, nil
}

func (s *TransactionService) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := config.DB.First(&transaction, id)
	if result.Error != nil {
		return nil, fmt.Errorf("transaction not found: %w", result.Error)
	}
	return &transaction, nil
}

func (s *TransactionService) Delete(id uint) error {
	result := config.DB.Delete(&models.Transaction{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete transaction: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (s *TransactionService) GetByType(transactionType string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := config.DB.Where("type = ?", transactionType).Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}
	return transactions, nil
}

func (s *TransactionService) GetByBank(bank string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := config.DB.Where("bank = ?", bank).Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}
	return transactions, nil
}

func (s *TransactionService) Update(id uint, updates map[string]interface{}) (*models.Transaction, error) {
	var transaction models.Transaction
	result := config.DB.First(&transaction, id)
	if result.Error != nil {
		return nil, fmt.Errorf("transaction not found: %w", result.Error)
	}

	result = config.DB.Model(&transaction).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", result.Error)
	}

	return &transaction, nil
}

type MonthlySummary struct {
	Month          string  `json:"month"`
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	NetAmount      float64 `json:"net_amount"`
	IncomeCount    int     `json:"income_count"`
	ExpenseCount   int     `json:"expense_count"`
	TransactionCount int   `json:"transaction_count"`
}

type CategorySummary struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

func (s *TransactionService) GetMonthlySummary(year int, month int) (*MonthlySummary, []CategorySummary, error) {
	// Build date pattern for SQLite (e.g., "%/11/2025" for November 2025)
	// Supports both DD/MM/YYYY and D/M/YYYY formats
	monthStr := fmt.Sprintf("%02d", month)
	yearStr := fmt.Sprintf("%d", year)

	// Match patterns like: 1/11/2025, 01/11/2025, 31/11/2025
	pattern := fmt.Sprintf("%%/%s/%s", monthStr, yearStr)

	var transactions []models.Transaction
	result := config.DB.Where("date LIKE ?", pattern).Find(&transactions)
	if result.Error != nil {
		return nil, nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}

	summary := &MonthlySummary{
		Month: fmt.Sprintf("%s/%s", monthStr, yearStr),
	}

	categoryMap := make(map[string]*CategorySummary)

	for _, t := range transactions {
		summary.TransactionCount++

		if t.Type == "income" {
			summary.TotalIncome += t.Amount
			summary.IncomeCount++
		} else if t.Type == "expense" {
			summary.TotalExpense += t.Amount
			summary.ExpenseCount++
		}

		// Category breakdown
		category := t.Category
		if category == "" {
			category = "ไม่มีหมวดหมู่"
		}

		key := category + "_" + t.Type
		if _, exists := categoryMap[key]; !exists {
			categoryMap[key] = &CategorySummary{
				Category: category,
				Type:     t.Type,
			}
		}
		categoryMap[key].Total += t.Amount
		categoryMap[key].Count++
	}

	summary.NetAmount = summary.TotalIncome - summary.TotalExpense

	// Convert map to slice
	var categories []CategorySummary
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}

	return summary, categories, nil
}
