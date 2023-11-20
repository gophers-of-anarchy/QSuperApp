package models

import "time"

type Airplane struct {
	ID        uint      `gorm:"primaryKey"`
	Type      uint      `gorm:"not null"`
	BasePrice float32   `gorm:"not null"`
	Number    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type CreateAirplaneDTO struct {
	Type      uint    `json:"type"`
	BasePrice float32 `gorm:"base_price"`
	Number    string  `gorm:"number"`
}

type UpdateAirplaneDTO struct {
	Type      uint    `json:"type"`
	BasePrice float32 `gorm:"base_price"`
	Number    string  `gorm:"number"`
	ID        uint    `json:"id"`
}

type GetAllAirplaneDTO struct {
	Airplanes  []Airplane
	TotalCount int64
}

const (
	Airliner = iota
	Trainer
	Military
	unknown
)

func CheckAirplaneType(airplaneType uint) uint {
	switch airplaneType {
	case 0:
		return Airliner
	case 1:
		return Trainer
	case 2:
		return Military
	default:
		return unknown
	}
}
