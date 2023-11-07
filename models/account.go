package models

import "time"

type Account struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Balance   float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	// Account belongs to a user
	User User `gorm:"foreignKey:UserID"`

	// Account has many cards
	Cards []Card `gorm:"foreignKey:AccountID"`
}

type AccountCreateRequest struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
