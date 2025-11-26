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
