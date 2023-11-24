package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	AirplaneID uint      `gorm:"not null"`
	Status     string    `gorm:"not null"`
	Number     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	//FromUser User `gorm:"foreignKey:UserID"`
	User     User     `gorm:"foreignKey:UserID"`
	Airplane Airplane `gorm:"foreignKey:AirplaneID"`
	// Order has many Payment
	Payment []Payment `gorm:"foreignKey:OrderID"`
}

type ResponseOrderStatus struct {
	Message     string
	OrderStatus string
}

type ResponseChangeOrderStatus struct {
	Message string
}
