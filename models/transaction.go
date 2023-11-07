package models

import "time"

type Transaction struct {
	ID         uint      `gorm:"primaryKey"`
	FromCardID uint      `gorm:"not null"`
	ToCardID   uint      `gorm:"not null"`
	Amount     float64   `gorm:"not null"`
	FeeID      uint      `gorm:"null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`

	// Transaction has a from_card
	FromCard Card `gorm:"foreignKey:FromCardID"`

	// Transaction has a to_card
	ToCard Card `gorm:"foreignKey:ToCardID"`

	// Transaction has one transaction fee
	Fee TransactionFee `gorm:"foreignKey:TransactionID"`
}

type TransferRequest struct {
	FromCardNumber string `json:"from_card_number"`
	ToCardNumber   string `json:"to_card_number"`
	Amount         string `json:"amount"`
}
