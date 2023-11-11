package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"QSuperApp/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ManageOrdersHandler(ctx echo.Context) error {
	var response models.ResponseOrderState
	orderIDString := ctx.Request().URL.Query().Get("order_id")
	status := ctx.Request().URL.Query().Get("status")
	orderID, convertErr := strconv.Atoi(orderIDString)
	if convertErr != nil {
		response.Message = messages.InternalError
		response.Status = "Rejected"
		ctx.JSON(http.StatusInternalServerError, response)
		return convertErr
	}
	var orders []models.Order
	result := configs.DB.Find(&orders)
	if result.Error != nil {
		response.Message = messages.InternalError
		response.Status = "Rejected"
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	var customerOrder models.Order
	for _, order := range orders {
		if int(order.ID) == orderID {
			configs.DB.Model(&order).UpdateColumn("status", status)
			customerOrder = order
			break
		}
	}
	var users []models.User
	result = configs.DB.Find(&users)
	if result.Error != nil {
		response.Message = messages.FailedToNotifyUser
		response.Status = status
		ctx.JSON(http.StatusInternalServerError, response)
		return result.Error
	}
	for _, user := range users {
		if int(user.ID) == customerOrder.UserID {
			subject := fmt.Sprintf("Your order %s", status)
			body := fmt.Sprintf("Your order with the following specifications was accepted by the admin:\nOrder_ID: %v\nAirplane_ID: %v\nPrice: %v\n", customerOrder.ID, customerOrder.AirplaneID, customerOrder.Price)
			sendErr := services.SendMail(subject, body, []string{user.Email})
			if sendErr != nil {
				response.Message = messages.FailedToSendEmail
				response.Status = status
				ctx.JSON(http.StatusInternalServerError, response)
				return result.Error
			}
			response.Message = messages.OrderStatusChangedSuccessfully
			response.Status = status
			ctx.JSON(http.StatusOK, response)
			return nil
		}
	}
	//if status == "Accepted" {
	//	for _, order := range orders {
	//		if int(order.ID) == orderID {
	//			configs.DB.Model(&order).UpdateColumn("status", "Approved")
	//		}
	//	}
	//} else {
	//	for _, order := range orders {
	//		if int(order.ID) == orderID {
	//			configs.DB.Model(&order).UpdateColumn("status", "Rejected")
	//		}
	//	}
	//}
	return nil
}
