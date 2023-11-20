package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"QSuperApp/services"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func DecideOrderStatusHandler(ctx echo.Context) error {
	var response models.ResponseOrderStatus
	orderIDString := ctx.Request().URL.Query().Get("order_id")
	status := ctx.Request().URL.Query().Get("status")
	orderID, convertErr := strconv.Atoi(orderIDString)
	if convertErr != nil {
		response.Message = messages.InternalError
		response.OrderStatus = "Rejected"
		ctx.JSON(http.StatusInternalServerError, response)
		return convertErr
	}
	if !(status == "Approved" || status == "Rejected") {
		response.Message = messages.WrongStatus
		ctx.JSON(http.StatusBadRequest, response)
		return errors.New(strings.ToLower(messages.WrongStatus))
	}
	var orders []models.Order
	result := configs.DB.Find(&orders)
	if result.Error != nil {
		response.Message = messages.InternalError
		response.OrderStatus = "Rejected"
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	exist := false
	var customerOrder models.Order
	for _, order := range orders {
		if int(order.ID) == orderID {
			configs.DB.Model(&order).UpdateColumn("status", status)
			customerOrder = order
			exist = true
			break
		}
	}
	if !exist {
		response.Message = messages.WrongOrderID
		response.OrderStatus = "Rejected"
		ctx.JSON(http.StatusBadRequest, response)
		return result.Error
	}
	var users []models.User
	result = configs.DB.Find(&users)
	if result.Error != nil {
		response.Message = messages.FailedToNotifyUser
		response.OrderStatus = status
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	var customs []models.Customization
	result = configs.DB.Find(&customs)
	if result.Error != nil {
		response.Message = messages.FailedToNotifyUser
		response.OrderStatus = status
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	var customerOrderCustom models.Customization
	for _, custom := range customs {
		if int(custom.OrderID) == orderID {
			customerOrderCustom = custom
			break
		}
	}
	for _, user := range users {
		if user.ID == customerOrder.UserID {
			subject := fmt.Sprintf("Your order %s", status)
			body := fmt.Sprintf("Your order with the following specifications was accepted by the admin:\nOrder-ID: %v\nAirplane-ID: %v\nVIP Seats: %v\nCockpit Facilities: %v\nAdditional Facilities: %v\n", customerOrder.ID, customerOrder.AirplaneID, customerOrderCustom.VIPSeatsCount, customerOrderCustom.CockpitFacilitiesLevel, customerOrderCustom.AdditionalFacilities)
			sendErr := services.SendMail(subject, body, []string{user.Email})
			if sendErr != nil {
				response.Message = messages.FailedToSendEmail
				response.OrderStatus = status
				ctx.JSON(http.StatusInternalServerError, response)
				return result.Error
			}
			response.Message = messages.OrderStatusChangedSuccessfully
			response.OrderStatus = status
			ctx.JSON(http.StatusOK, response)
			return nil
		}
	}
	return nil
}

func ChangeOrderStatusHandler(ctx echo.Context) error {
	var response models.ResponseChangeOrderStatus
	orderIDString := ctx.Request().URL.Query().Get("order_id")
	status := ctx.Request().URL.Query().Get("status")
	orderID, convertErr := strconv.Atoi(orderIDString)
	if convertErr != nil {
		response.Message = messages.InternalError
		ctx.JSON(http.StatusInternalServerError, response)
		return convertErr
	}
	if !(status == "Under construction" || status == "Built! ready for delivery") {
		response.Message = messages.WrongStatus
		ctx.JSON(http.StatusBadRequest, response)
		return errors.New(strings.ToLower(messages.WrongStatus))
	}
	var orders []models.Order
	result := configs.DB.Find(&orders)
	if result.Error != nil {
		response.Message = messages.InternalError
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	exist := false
	var customerOrder models.Order
	for _, order := range orders {
		if int(order.ID) == orderID {
			if order.Status == "Approved" {
				if status != "Under construction" {
					response.Message = messages.FailedToChangeStatus
					ctx.JSON(http.StatusBadRequest, response)
					return result.Error
				}
			} else if order.Status == "Under construction" {
				if status != "Built! ready for delivery" {
					response.Message = messages.FailedToChangeStatus
					ctx.JSON(http.StatusBadRequest, response)
					return result.Error
				}
			}
			configs.DB.Model(&order).UpdateColumn("status", status)
			customerOrder = order
			exist = true
			break
		}
	}
	if !exist {
		response.Message = messages.WrongOrderID
		ctx.JSON(http.StatusBadRequest, response)
		return result.Error
	}
	var users []models.User
	result = configs.DB.Find(&users)
	if result.Error != nil {
		response.Message = messages.FailedToNotifyUser
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	var customs []models.Customization
	result = configs.DB.Find(&customs)
	if result.Error != nil {
		response.Message = messages.FailedToNotifyUser
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	var customerOrderCustom models.Customization
	for _, custom := range customs {
		if int(custom.OrderID) == orderID {
			customerOrderCustom = custom
			break
		}
	}
	for _, user := range users {
		if user.ID == customerOrder.UserID {
			if status == "Under construction" {
				subject := fmt.Sprintf("Your order is under construction")
				body := fmt.Sprintf("Your order with the following specifications is getting ready:\nOrder-ID: %v\nAirplane-ID: %v\nVIP Seats: %v\nCockpit Facilities: %v\nAdditional Facilities: %v\n", customerOrder.ID, customerOrder.AirplaneID, customerOrderCustom.VIPSeatsCount, customerOrderCustom.CockpitFacilitiesLevel, customerOrderCustom.AdditionalFacilities)
				sendErr := services.SendMail(subject, body, []string{user.Email})
				if sendErr != nil {
					response.Message = messages.FailedToSendEmail
					ctx.JSON(http.StatusInternalServerError, response)
					return result.Error
				}
				response.Message = messages.OrderStatusChangedSuccessfully
				ctx.JSON(http.StatusOK, response)
				return nil
			} else {
				subject := fmt.Sprintf("Your order has been completed")
				body := fmt.Sprintf("Your order with the following specifications has been completed and is ready for delivery:\nOrder-ID: %v\nAirplane-ID: %v\nVIP Seats: %v\nCockpit Facilities: %v\nAdditional Facilities: %v\n", customerOrder.ID, customerOrder.AirplaneID, customerOrderCustom.VIPSeatsCount, customerOrderCustom.CockpitFacilitiesLevel, customerOrderCustom.AdditionalFacilities)
				sendErr := services.SendMail(subject, body, []string{user.Email})
				if sendErr != nil {
					response.Message = messages.FailedToSendEmail
					ctx.JSON(http.StatusInternalServerError, response)
					return result.Error
				}
				response.Message = messages.OrderStatusChangedSuccessfully
				ctx.JSON(http.StatusOK, response)
				return nil
			}
		}
	}
	return nil
}
