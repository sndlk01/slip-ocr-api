package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Name            string         `gorm:"type:varchar(200);not null" json:"name"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Category        string         `gorm:"type:varchar(100)" json:"category"`
	BillingCycle    string         `gorm:"type:varchar(20);not null" json:"billing_cycle"` // monthly, yearly
	NextBillingDate string         `gorm:"type:varchar(20)" json:"next_billing_date"`      // DD/MM/YYYY
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	AutoDetected    bool           `gorm:"default:false" json:"auto_detected"` // ถูก detect จาก OCR มั้ย
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
