package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Add(ctx echo.Context) error {
	airplaneDTO := models.CreateAirplaneDTO{}
	err := ctx.Bind(&airplaneDTO)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	airplane := models.Airplane{
		Type:      models.CheckAirplaneType(airplaneDTO.Type),
		BasePrice: airplaneDTO.BasePrice,
		Number:    airplaneDTO.Number,
	}
	if result := configs.DB.Save(&airplane); result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, &airplane)
}

func Update(ctx echo.Context) error {
	airplaneDTO := models.UpdateAirplaneDTO{}
	err := ctx.Bind(&airplaneDTO)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	airplane := models.Airplane{
		Type:      models.CheckAirplaneType(airplaneDTO.Type),
		BasePrice: airplaneDTO.BasePrice,
		Number:    airplaneDTO.Number,
	}
	if err := configs.DB.First(&models.Airplane{ID: airplaneDTO.ID}).Updates(&airplane).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, &airplaneDTO)
}

func GetAll(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParams().Get("page"))
	pageSize, _ := strconv.Atoi(ctx.QueryParams().Get("pageSize"))
	sort := ctx.QueryParams().Get("sortBy")
	order := ctx.QueryParams().Get("sortOrder")

	var airplanes []models.Airplane
	count := configs.DB.Find(&airplanes).RowsAffected
	configs.DB.Order(sort + " " + order).Limit(pageSize).Offset(page * pageSize).Find(&airplanes)

	return ctx.JSON(
		http.StatusOK,
		models.GetAllAirplaneDTO{
			Airplanes:  airplanes,
			TotalCount: count,
		},
	)
}

func Get(ctx echo.Context) error {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	airplane := models.Airplane{ID: uint(id)}
	result := configs.DB.First(&airplane)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return ctx.JSON(http.StatusOK, &airplane)
}

func Delete(ctx echo.Context) error {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	airplane := models.Airplane{ID: uint(id)}
	if err := configs.DB.Delete(&airplane).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, nil)
}
