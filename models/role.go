package models

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	RoleTitle string `gorm:"not null"`
	// Role has many RoleUser
	RoleUser []RoleUser `gorm:"foreignKey:RoleID"`
}
