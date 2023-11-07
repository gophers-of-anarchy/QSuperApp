package tests

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTransferHandler(t *testing.T) {
	// Connect to the database
	configs.ConnectToDatabase()

	// Connect redis
	configs.ConnectToRedis()

	// Create a new user for testing
	user := models.User{
		Username:  "testuser",
		Password:  "testpassword",
		Email:     "testuser@example.com",
		Cellphone: "09123456789",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&user)

	// Create two new accounts for testing
	account1 := models.Account{
		UserID:    user.ID,
		Name:      "Test Account",
		Balance:   15000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&account1)

	account2 := models.Account{
		UserID:    user.ID,
		Name:      "Test Account",
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&account2)

	// Create two new cards for testing
	card1 := models.Card{
		AccountID: account1.ID,
		Number:    "6037991161333780",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&card1)

	card2 := models.Card{
		AccountID: account2.ID,
		Number:    "6037991161333781",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&card2)

	// Create a new transfer request for testing
	transferReq := models.TransferRequest{
		FromCardNumber: card1.Number,
		ToCardNumber:   card2.Number,
		Amount:         "5000",
	}

	// Encode the transfer request as JSON
	reqBody, _ := json.Marshal(transferReq)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/transfer", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the user_id value in the request context
	ctx := context.WithValue(req.Context(), "user_id", user.ID)
	req = req.WithContext(ctx)

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the TransferHandler function to handle the request
	http.HandlerFunc(controllers.TransferHandler).ServeHTTP(rr, req)

	// Assert that the HTTP response has a status code of 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the response body contains the expected message
	expectedRes := messages.TransferSuccessful
	assert.Equal(t, expectedRes, rr.Body.String())

	// Query the database to assert that the transfer has been processed
	var fromCard models.Card
	configs.DB.Preload("Account").Preload("OutgoingTransactions").First(&fromCard, "number = ?", transferReq.FromCardNumber)
	assert.NotNil(t, fromCard)
	assert.Equal(t, float64(9500), fromCard.Account.Balance)
	assert.Len(t, fromCard.OutgoingTransactions, 1)

	var toCard models.Card
	configs.DB.Preload("Account").Preload("IncomingTransactions").First(&toCard, "number = ?", transferReq.ToCardNumber)
	assert.NotNil(t, toCard)
	assert.Equal(t, float64(5000), toCard.Account.Balance)
	assert.Len(t, toCard.IncomingTransactions, 1)
}
