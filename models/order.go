package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Status    string    `gorm:"not null"`
	Number    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	//FromUser User `gorm:"foreignKey:UserID"`
	User User `gorm:"foreignKey:UserID"`

	// Order has many Payment
	Payment []Payment `gorm:"foreignKey:OrderID"`
}
