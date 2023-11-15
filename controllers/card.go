package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/helpers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"encoding/json"
	"net/http"
	"time"
)

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the card details
	var cardCreateReq models.CardCreateRequest
	err := json.NewDecoder(r.Body).Decode(&cardCreateReq)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Check if the account exists
	var account models.Account
	if result := configs.DB.First(&account, cardCreateReq.AccountID); result.Error != nil {
		http.Error(w, messages.AccountNotFound, http.StatusNotFound)
		return
	}

	var card models.Card

	// Set the account ID for the new card
	card.AccountID = account.ID

	// Clean card number
	cardNumber := helpers.CardNumberToEnglish(cardCreateReq.Number)

	// Validate card number
	if !helpers.CardIsValid(cardNumber) {
		http.Error(w, messages.InvalidCardNumber, http.StatusBadRequest)
		return
	}

	card.Number = cardNumber

	// Set the creation time
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	// Save the new card to the database
	if result := configs.DB.Create(&card); result.Error != nil {
		http.Error(w, messages.FailedToCreateCard, http.StatusInternalServerError)
		return
	}

	// Respond with message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(messages.CardCreatedSuccessfully))
}
