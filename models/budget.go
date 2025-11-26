package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Category     string         `gorm:"type:varchar(100);not null" json:"category"`
	MonthlyLimit float64        `gorm:"not null" json:"monthly_limit"`
	Month        int            `gorm:"not null" json:"month"` // 1-12
	Year         int            `gorm:"not null" json:"year"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Budget) TableName() string {
	return "budgets"
}
