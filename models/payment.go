package models

import "time"

type Payment struct {
	ID            uint      `gorm:"primaryKey"`
	Amount        float64   `gorm:"not null"`
	PaymentType   int64     `gorm:"not null"`
	PaymentStatus int64     `gorm:"not null"`
	OrderID       uint      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	// Payment belongs to an Order
	Order Order `gorm:"foreignKey:OrderID"`
}
