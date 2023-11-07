package models

import "time"

type Card struct {
	ID        uint      `gorm:"primaryKey"`
	AccountID uint      `gorm:"not null"`
	Number    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	// Card belongs to an account
	Account Account `gorm:"foreignKey:AccountID"`

	// Card has many transactions as from_card
	OutgoingTransactions []Transaction `gorm:"foreignKey:FromCardID"`

	// Card has many transactions as to_card
	IncomingTransactions []Transaction `gorm:"foreignKey:ToCardID"`
}

type CardCreateRequest struct {
	AccountID string `json:"account_id"`
	Number    string `json:"card_number"`
}
