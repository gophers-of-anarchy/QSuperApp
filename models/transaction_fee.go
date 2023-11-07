package models

import "time"

type TransactionFee struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID uint      `gorm:"not null"`
	Amount        float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}
