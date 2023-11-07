package tests

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateCardHandler(t *testing.T) {
	// Connect postgres
	configs.ConnectToDatabase()

	// Create a new account for testing
	account := models.Account{
		UserID:    1,
		Name:      "Test Account",
		Balance:   1000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	configs.DB.Create(&account)

	// Create a new card creation request for testing
	cardCreateReq := models.CardCreateRequest{
		AccountID: strconv.Itoa(int(account.ID)),
		Number:    "۶۰۳۷-۹۹۱۱-۶۱۳۳-۳۷۸۰",
	}

	// Encode the card creation request as JSON
	reqBody, _ := json.Marshal(cardCreateReq)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/create_card", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the CreateCardHandler function to handle the request
	http.HandlerFunc(controllers.CreateCardHandler).ServeHTTP(rr, req)

	// Assert that the HTTP response has a status code of 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the response body contains the expected message
	expectedRes := messages.CardCreatedSuccessfully
	assert.Equal(t, expectedRes, rr.Body.String())

	// Query the database to assert that the new card has been created
	var card models.Card
	configs.DB.First(&card, "number = ?", "6037991161333780")
	assert.NotNil(t, card)
}
