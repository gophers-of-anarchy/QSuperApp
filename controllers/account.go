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

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value("user_id").(uint)
	if userID == 0 {
		http.Error(w, messages.Unauthorized, http.StatusUnauthorized)
		return
	}
	// Parse the request body to get the account details
	var accountCreateReq models.AccountCreateRequest
	err := json.NewDecoder(r.Body).Decode(&accountCreateReq)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Create a new account
	account := models.Account{
		UserID:    userID,
		Name:      accountCreateReq.Name,
		Balance:   accountCreateReq.Balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the new account to the database
	if result := configs.DB.Create(&account); result.Error != nil {
		http.Error(w, messages.FailedToCreateAccount, http.StatusInternalServerError)
		return
	}

	// Respond with the new account
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(messages.AccountCreatedSuccessfully, account.ID)))
}
