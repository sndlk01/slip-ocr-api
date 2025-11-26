package services

import (
	"fmt"
	"ocr-api/config"
	"ocr-api/models"
	"regexp"
	"strings"
)

type SubscriptionService struct{}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}

// Subscription patterns for auto-detection from OCR
var subscriptionPatterns = []struct {
	Name     string
	Patterns []string
	Category string
}{
	{
		Name: "Netflix",
		Patterns: []string{
			"(?i)netflix",
		},
		Category: "บันเทิง",
	},
	{
		Name: "Spotify",
		Patterns: []string{
			"(?i)spotify",
		},
		Category: "บันเทิง",
	},
	{
		Name: "YouTube Premium",
		Patterns: []string{
			"(?i)youtube\\s*premium",
			"(?i)youtube\\s*music",
		},
		Category: "บันเทิง",
	},
	{
		Name: "LINE MAN",
		Patterns: []string{
			"(?i)line\\s*man",
			"(?i)lineman",
		},
		Category: "สมาชิก",
	},
	{
		Name: "Grab Unlimited",
		Patterns: []string{
			"(?i)grab\\s*unlimited",
		},
		Category: "สมาชิก",
	},
	{
		Name: "True ID",
		Patterns: []string{
			"(?i)true\\s*id",
			"(?i)trueid",
		},
		Category: "บันเทิง",
	},
	{
		Name: "Disney+",
		Patterns: []string{
			"(?i)disney\\s*\\+",
			"(?i)disney\\s*plus",
		},
		Category: "บันเทิง",
	},
	{
		Name: "iCloud",
		Patterns: []string{
			"(?i)icloud",
			"(?i)apple\\s*storage",
		},
		Category: "คลาวด์",
	},
	{
		Name: "Google One",
		Patterns: []string{
			"(?i)google\\s*one",
		},
		Category: "คลาวด์",
	},
	{
		Name: "Adobe",
		Patterns: []string{
			"(?i)adobe",
			"(?i)photoshop",
		},
		Category: "ซอฟต์แวร์",
	},
}

// DetectSubscription tries to detect subscription service from OCR text
func (s *SubscriptionService) DetectSubscription(ocrText string, amount float64) *models.Subscription {
	for _, sp := range subscriptionPatterns {
		for _, pattern := range sp.Patterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(ocrText) {
				return &models.Subscription{
					Name:         sp.Name,
					Amount:       amount,
					Category:     sp.Category,
					BillingCycle: "monthly", // default
					IsActive:     true,
					AutoDetected: true,
				}
			}
		}
	}
	return nil
}

func (s *SubscriptionService) Create(subscription *models.Subscription) error {
	result := config.DB.Create(subscription)
	if result.Error != nil {
		return fmt.Errorf("failed to create subscription: %w", result.Error)
	}
	return nil
}

func (s *SubscriptionService) GetAll() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := config.DB.Order("is_active DESC, name ASC").Find(&subscriptions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", result.Error)
	}
	return subscriptions, nil
}

func (s *SubscriptionService) GetActive() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := config.DB.Where("is_active = ?", true).
		Order("name ASC").Find(&subscriptions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", result.Error)
	}
	return subscriptions, nil
}

func (s *SubscriptionService) GetByID(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	result := config.DB.First(&subscription, id)
	if result.Error != nil {
		return nil, fmt.Errorf("subscription not found: %w", result.Error)
	}
	return &subscription, nil
}

func (s *SubscriptionService) Update(id uint, updates map[string]interface{}) (*models.Subscription, error) {
	var subscription models.Subscription
	result := config.DB.First(&subscription, id)
	if result.Error != nil {
		return nil, fmt.Errorf("subscription not found: %w", result.Error)
	}

	result = config.DB.Model(&subscription).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", result.Error)
	}

	return &subscription, nil
}

func (s *SubscriptionService) Delete(id uint) error {
	result := config.DB.Delete(&models.Subscription{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// GetUpcoming returns subscriptions with billing dates in the next N days
func (s *SubscriptionService) GetUpcoming(days int) ([]models.Subscription, error) {
	// This is a simplified version - you might want to parse dates properly
	var subscriptions []models.Subscription
	result := config.DB.Where("is_active = ? AND next_billing_date != ''", true).
		Order("next_billing_date ASC").Find(&subscriptions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get upcoming subscriptions: %w", result.Error)
	}
	return subscriptions, nil
}

// CalculateMonthlyTotal calculates total monthly subscription cost
func (s *SubscriptionService) CalculateMonthlyTotal() (float64, error) {
	var total float64
	result := config.DB.Model(&models.Subscription{}).
		Where("is_active = ? AND billing_cycle = ?", true, "monthly").
		Select("COALESCE(SUM(amount), 0)").Scan(&total)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to calculate total: %w", result.Error)
	}
	return total, nil
}

// SuggestCategory suggests a category based on receiver name from OCR
func SuggestSubscriptionCategory(receiver string) string {
	receiver = strings.ToLower(receiver)

	if strings.Contains(receiver, "netflix") || strings.Contains(receiver, "spotify") ||
		strings.Contains(receiver, "youtube") || strings.Contains(receiver, "true") ||
		strings.Contains(receiver, "disney") {
		return "บันเทิง"
	}

	if strings.Contains(receiver, "icloud") || strings.Contains(receiver, "google") ||
		strings.Contains(receiver, "dropbox") {
		return "คลาวด์"
	}

	if strings.Contains(receiver, "adobe") || strings.Contains(receiver, "microsoft") {
		return "ซอฟต์แวร์"
	}

	return "สมาชิก"
}
