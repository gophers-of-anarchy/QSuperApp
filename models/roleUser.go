package models

type RoleUser struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	RoleID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`
	Role   Role `gorm:"foreignKey:RoleID"`
}
