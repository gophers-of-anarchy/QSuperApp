package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/helpers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"QSuperApp/services"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AdvancePaymentHandler(ctx echo.Context) error {
	// Get user ID from context
	userID := ctx.Get("user_id").(string)

	var req models.AdvancePaymentRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	var order models.Order
	orderResult := configs.DB.Where("id = ? AND user_id = ?", req.OrderId, userID).Preload("User").First(&order)
	if errors.Is(orderResult.Error, gorm.ErrRecordNotFound) {
		log.Println("Order not found:", orderResult.Error)
		return ctx.JSON(http.StatusNotFound, messages.OrderNotFound)
	}

	var customization models.Customization
	customizationResult := configs.DB.Where("order_id = ?", order.ID).First(&customization)
	if errors.Is(customizationResult.Error, gorm.ErrRecordNotFound) {
		log.Println("Customization not found:", customizationResult.Error)
		return ctx.JSON(http.StatusNotFound, messages.OrderNotFound)
	}

	totalPrice := helpers.CalculateAirplaneTotalPrice(customization)
	advanceAmount := totalPrice * 0.5

	// create a payment record
	payment := models.Payment{
		Amount:        advanceAmount,
		PaymentType:   models.AdvancePayment,
		PaymentStatus: models.PaymentPending,
		OrderID:       order.ID,
		CreatedAt:     time.Now(),
		Order:         order,
	}
	if result := configs.DB.Save(&payment); result.Error != nil {
		log.Println("Error saving payment:", result.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	paymentResp, err := services.Payment(advanceAmount, payment.ID)
	if err != nil {
		log.Println("Error making payment:", err)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	if paymentResp.ResCode != 0 {
		if result := configs.DB.Model(&payment).Updates(models.Payment{PaymentStatus: models.PaymentFailed, UpdatedAt: time.Now()}); result.Error != nil {
			log.Println("Error updating payment:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
		}
		log.Println("Error making payment:", paymentResp.Description)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	return ctx.JSON(http.StatusOK, paymentResp)
}

func FinalPaymentHandler(ctx echo.Context) error {
	// Get user ID from context
	userID := ctx.Get("user_id").(string)

	var req models.AdvancePaymentRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	var order models.Order
	orderResult := configs.DB.Where("id = ? AND user_id = ?", req.OrderId, userID).Preload("User").First(&order)
	if errors.Is(orderResult.Error, gorm.ErrRecordNotFound) {
		log.Println("Order not found:", orderResult.Error)
		return ctx.JSON(http.StatusNotFound, messages.OrderNotFound)
	}

	var customization models.Customization
	customizationResult := configs.DB.Where("order_id = ?", order.ID).First(&customization)
	if errors.Is(customizationResult.Error, gorm.ErrRecordNotFound) {
		log.Println("Customization not found:", customizationResult.Error)
		return ctx.JSON(http.StatusNotFound, messages.OrderNotFound)
	}

	totalPrice := helpers.CalculateAirplaneTotalPrice(customization)
	var paymentAmount float64
	var advancePayment models.Payment
	advancePaymentResult := configs.DB.Where("order_id = ? AND payment_type = ?", order.ID, models.AdvancePayment).First(&advancePayment)
	if errors.Is(advancePaymentResult.Error, gorm.ErrRecordNotFound) {
		log.Println("Advance payment not found:", advancePaymentResult.Error)
		return ctx.JSON(http.StatusNotFound, messages.PaymentNotFound)
	}
	if advancePayment.PaymentStatus != models.PaymentCompleted {
		paymentAmount = totalPrice
	}else {
		paymentAmount = totalPrice * 0.5
	}

	// create a payment record
	payment := models.Payment{
		Amount:        paymentAmount,
		PaymentType:   models.FinalPayment,
		PaymentStatus: models.PaymentPending,
		OrderID:       order.ID,
		CreatedAt:     time.Now(),
		Order:         order,
	}
	if result := configs.DB.Save(&payment); result.Error != nil {
		log.Println("Error saving payment:", result.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	paymentResp, err := services.Payment(paymentAmount, payment.ID)
	if err != nil {
		log.Println("Error making payment:", err)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	if paymentResp.ResCode != 0 {
		if result := configs.DB.Model(&payment).Updates(models.Payment{PaymentStatus: models.PaymentFailed, UpdatedAt: time.Now()}); result.Error != nil {
			log.Println("Error updating payment:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
		}
		log.Println("Error making payment:", paymentResp.Description)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	return ctx.JSON(http.StatusOK, paymentResp)
}

func VerifyPaymentPageHandler(ctx echo.Context) error {
	requestData := make(map[string]interface{})

	if err := ctx.Bind(&requestData); err != nil {
		return err
	}

	var payment models.Payment
	paymentResult := configs.DB.Where("id = ?", requestData["OrderId"]).First(&payment)
	if paymentResult.Error != nil {
		log.Println("Error finding payment:", paymentResult.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentNotFound)
	}

	if requestData["ResCode"] != "0" {
		if result := configs.DB.Model(&payment).Updates(models.Payment{PaymentStatus: models.PaymentFailed, UpdatedAt: time.Now()}); result.Error != nil {
			log.Println("Error updating payment:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
		}
		return ctx.Render(http.StatusBadRequest, "failed_payment.html", requestData)
	}

	data := map[string]interface{}{
		"Title":   "Successfull Payment",
		"OrderID": requestData["OrderId"],
		"Amount":  payment.Amount,
		"Token":   requestData["Token"],
	}
	return ctx.Render(http.StatusOK, "verify.html", data)
}

func VerifyPaymentHandler(ctx echo.Context) error {
	requestData := make(map[string]interface{})

	if err := ctx.Bind(&requestData); err != nil {
		return err
	}

	var payment models.Payment
	paymentResult := configs.DB.Where("id = ?", requestData["OrderId"]).First(&payment)
	if paymentResult.Error != nil {
		log.Println("Error finding payment:", paymentResult.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentNotFound)
	}

	resp, err := services.VerifyPayment(requestData["Token"].(string))
	if err != nil {
		log.Println("Error verifying payment:", err)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentVerificationError)
	}

	if resp.ResCode != 0 {
		if result := configs.DB.Model(&payment).Updates(models.Payment{PaymentStatus: models.PaymentFailed, UpdatedAt: time.Now()}); result.Error != nil {
			log.Println("Error updating payment:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
		}
		log.Println("Error verifying payment:", resp.Description)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentVerificationError)
	}
	if result := configs.DB.Model(&payment).Updates(models.Payment{PaymentStatus: models.PaymentCompleted, UpdatedAt: time.Now()}); result.Error != nil {
		log.Println("Error updating payment:", result.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentError)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func GetOrderPaymentDetailsHandler(ctx echo.Context) error {
	orderId := ctx.Param("orderId")

	var payments []models.Payment
	paymentResult := configs.DB.Where("order_id = ?", orderId).Find(&payments)
	if paymentResult.Error != nil {
		log.Println("Error finding payments:", paymentResult.Error)
		return ctx.JSON(http.StatusInternalServerError, messages.PaymentNotFound)
	}
	// return the payments as result
	return ctx.JSON(http.StatusOK, payments)
}
