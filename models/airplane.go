package models

import "time"

type Airplane struct {
	ID uint `gorm:"primaryKey"`

	Type      uint      `gorm:"not null"`
	BasePrice float32   `gorm:"not null"`
	Number    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
