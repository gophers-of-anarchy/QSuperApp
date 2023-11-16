package models

import "time"

type Authentication struct {
	TokenID uint `gorm:"primaryKey"`
	UserID  uint `gorm:"not null"`

	JWTToken  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	Expiry    time.Time `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	//FromUser User `gorm:"foreignKey:FromUserID"`
}
