package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/helpers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := strconv.Itoa(int(r.Context().Value("user_id").(uint)))
	if userID == "" {
		http.Error(w, messages.Unauthorized, http.StatusUnauthorized)
		return
	}
	// Implement rate limiting
	key := fmt.Sprintf(os.Getenv("CACHE_KEY_USER_RATE_LIMIT"), userID)
	currentCount, err := configs.Redis.Incr(r.Context(), key).Result()
	if err != nil {
		http.Error(w, messages.InternalError, http.StatusInternalServerError)
		return
	}
	if currentCount > 10 {
		http.Error(w, messages.TooManyRequests, http.StatusTooManyRequests)
		return
	}
	if currentCount == 1 {
		_, err := configs.Redis.Expire(r.Context(), key, time.Minute).Result()
		if err != nil {
			http.Error(w, messages.InternalError, http.StatusInternalServerError)
			return
		}
	}

	// Parse request body
	var transferRequest models.TransferRequest
	err = json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Get the sender and receiver cards
	var fromCard, toCard models.Card
	configs.DB.Where("number = ?", transferRequest.FromCardNumber).Preload("Account").First(&fromCard)
	configs.DB.Where("number = ?", transferRequest.ToCardNumber).Preload("Account").First(&toCard)

	// Check if both cards exist
	if fromCard.ID == 0 || toCard.ID == 0 {
		http.Error(w, messages.InvalidCardNumber, http.StatusBadRequest)
		return
	}

	// Parse transfer amount
	transferAmount, err := helpers.AmountToEnglish(transferRequest.Amount)
	if err != nil {
		http.Error(w, messages.InvalidTransferAmount, http.StatusBadRequest)
		return
	}

	// Check if the sender has enough balance
	if fromCard.Account.Balance < transferAmount+500 {
		http.Error(w, messages.InsufficientBalance, http.StatusBadRequest)
		return
	}

	// Check if the transfer amount is within the limits
	if transferAmount < 1000 || transferAmount > 500000 {
		http.Error(w, messages.InvalidTransferAmount, http.StatusBadRequest)
		return
	}

	// Start a database transaction
	tx := configs.DB.Begin()

	// Deduct the transfer amount and fee from the sender's card
	fromCard.Account.Balance -= transferAmount + 500
	if err := tx.Save(&fromCard.Account).Error; err != nil {
		tx.Rollback()
		http.Error(w, messages.FailedToUpdateSender, http.StatusInternalServerError)
		return
	}

	// Add the transfer amount to the receiver's card
	toCard.Account.Balance += transferAmount
	if err := tx.Save(&toCard.Account).Error; err != nil {
		tx.Rollback()
		http.Error(w, messages.FailedToUpdateReceiver, http.StatusInternalServerError)
		return
	}

	// Add a new transaction record
	transaction := models.Transaction{
		FromCardID: fromCard.ID,
		ToCardID:   toCard.ID,
		Amount:     transferAmount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		http.Error(w, messages.FailedToCreateTransaction, http.StatusInternalServerError)
		return
	}

	// Add a transaction fee record
	fee := models.TransactionFee{
		TransactionID: transaction.ID,
		Amount:        500,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&fee).Error; err != nil {
		tx.Rollback()
		http.Error(w, messages.FailedToCreateTransactionFee, http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		http.Error(w, messages.FailedToCommitTransaction, http.StatusInternalServerError)
		return
	}

	//Send SMS notifications to sender and receiver
	//msg := fmt.Sprintf(messages.TransactionSenderSMSText, fromCard.Number, transferAmount, fromCard.Account.Balance)
	//err = services.SMS.SendSMS(fromCard.Account.User.Cellphone, msg)
	//if err != nil {
	//	http.Error(w, messages.FailedToSendSenderSMS, http.StatusInternalServerError)
	//	return
	//}
	//
	//msg = fmt.Sprintf(messages.TransactionReceiverSMSText, toCard.Number, transferAmount, toCard.Account.Balance)
	//err = services.SMS.SendSMS(toCard.Account.User.Cellphone, msg)
	//if err != nil {
	//	http.Error(w, messages.FailedToSendReceiverSMS, http.StatusInternalServerError)
	//	return
	//}

	// Add the transaction to the user's transaction set
	transactionJSON, _ := json.Marshal(transaction)
	cacheKey := fmt.Sprintf(os.Getenv("CACHE_KEY_USER_TRANSACTIONS"), userID)
	configs.Redis.LPush(r.Context(), cacheKey, string(transactionJSON))

	// Increment the user's transaction count in the last 10 minutes
	configs.Redis.ZIncrBy(r.Context(), os.Getenv("CACHE_KEY_TRANSACTIONS_COUNT"), 1, userID)

	// Respond with message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(messages.TransferSuccessful))
}

func TopUsersHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	tenMinutesAgo := now.Add(-10 * time.Minute)

	var users []models.User
	err := configs.DB.
		Select("users.*, COUNT(transactions.id) as transaction_count").
		Joins("JOIN accounts ON users.id = accounts.user_id").
		Joins("JOIN cards ON accounts.id = cards.account_id").
		Joins("JOIN transactions ON cards.id = transactions.from_card_id OR cards.id = transactions.to_card_id").
		Where("transactions.created_at BETWEEN ? AND ?", tenMinutesAgo, now).
		Group("users.id").
		Order("transaction_count DESC").
		Limit(3).
		Find(&users).
		Error
	if err != nil {
		http.Error(w, messages.InternalError, http.StatusInternalServerError)
		return
	}

	type userWithTransactions struct {
		User         models.User          `json:"user"`
		Transactions []models.Transaction `json:"transactions"`
	}

	var usersWithTransactions []userWithTransactions
	for _, user := range users {
		var transactions []models.Transaction
		err := configs.DB.
			Preload("FromCard.Account.User").
			Preload("ToCard.Account.User").
			Where("transactions.created_at BETWEEN ? AND ?", tenMinutesAgo, now).
			Where("from_card_id IN (SELECT id FROM cards WHERE account_id IN (SELECT id FROM accounts WHERE user_id = ?)) OR to_card_id IN (SELECT id FROM cards WHERE account_id IN (SELECT id FROM accounts WHERE user_id = ?))", user.ID, user.ID).
			Order("created_at DESC").
			Limit(10).
			Find(&transactions).
			Error
		if err != nil {
			http.Error(w, messages.InternalError, http.StatusInternalServerError)
			return
		}
		usersWithTransactions = append(usersWithTransactions, userWithTransactions{User: user, Transactions: transactions})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersWithTransactions)
}
