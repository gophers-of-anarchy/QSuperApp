package models

import "time"

type Customization struct {
	ID                             uint      `gorm:"primaryKey"`
	OrderID                        uint      `gorm:"not null"`
	VIPSeatsCount                  int       `gorm:"not null"`
	ExteriorPaintingDesignAndColor int       `gorm:"null"`
	SeatConfiguration              int       `gorm:"not null"`
	AdditionalFacilities           int       `gorm:"not null"`
	CockpitFacilitiesLeve          int       `gorm:"not null"`
	CreatedAt                      time.Time `gorm:"not null"`
	UpdatedAt                      time.Time `gorm:"not null"`
	//FromUser Order `gorm:"foreignKey:OrderID"`
	Order Order `gorm:"foreignKey:OrderID"`
}
