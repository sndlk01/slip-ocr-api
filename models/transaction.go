package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Type       string         `gorm:"type:varchar(10);not null" json:"type"`
	Amount     float64        `gorm:"not null" json:"amount"`
	Date       string         `gorm:"type:varchar(20)" json:"date"`
	Time       string         `gorm:"type:varchar(20)" json:"time,omitempty"`
	Reference  string         `gorm:"type:varchar(100)" json:"reference,omitempty"`
	Bank       string         `gorm:"type:varchar(50)" json:"bank,omitempty"`
	Sender     string         `gorm:"type:varchar(200)" json:"sender,omitempty"`
	Receiver   string         `gorm:"type:varchar(200)" json:"receiver,omitempty"`
	Category   string         `gorm:"type:varchar(100)" json:"category"`
	Detail     string         `gorm:"type:text" json:"detail"`
	RawOCRText string         `gorm:"type:text" json:"raw_ocr_text,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Transaction) TableName() string {
	return "transactions"
}
