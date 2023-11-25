package helpers

import (
	"QSuperApp/models"
)

const (
	VIPSeat                        = 10000
	ExteriorPaintingDesignAndColor = 20000
	SeatConfiguration              = 30000
	AdditionalFacilities           = 40000
	CockpitFacilitiesLevel         = 50000
	BasePrice                      = 1000000.0
)

func CalculateAirplaneTotalPrice(customization models.Customization) float64 {
	var totalPrice = BasePrice
	totalPrice += float64(customization.VIPSeatsCount) * VIPSeat
	totalPrice += float64(customization.ExteriorPaintingDesignAndColor) * ExteriorPaintingDesignAndColor
	totalPrice += float64(customization.SeatConfiguration) * SeatConfiguration
	totalPrice += float64(customization.AdditionalFacilities) * AdditionalFacilities
	totalPrice += float64(customization.CockpitFacilitiesLevel) * CockpitFacilitiesLevel
	return totalPrice
}
