package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value("user_id").(uint)
	if userID == 0 {
		http.Error(w, messages.Unauthorized, http.StatusUnauthorized)
		return
	}

	// Decode the JSON request body into OrderCreateRequest
	var orderCreateReq models.OrderCreateRequest
	err := json.NewDecoder(r.Body).Decode(&orderCreateReq)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := orderCreateReq.Validate(configs.DB); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new order
	order := models.Order{
		UserID:     userID,
		AirplaneID: orderCreateReq.AirplaneID,
		Number:     generateOrderNumber(), // You may need to implement a function to generate order numbers
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save the new order to the database
	if result := configs.DB.Create(&order); result.Error != nil {
		http.Error(w, messages.FailedToCreateOrder, http.StatusInternalServerError)
		return
	}

	// Respond with the new order ID
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(messages.OrderCreatedSuccessfully, order.ID)))
}

// generateOrderNumber is a placeholder function, you may need to implement your logic to generate order numbers.
func generateOrderNumber() string {
	// Implement your logic here
	return "ORD-" + time.Now().Format("20060102-150405")
}
