package tests

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"QSuperApp/messages"
	"QSuperApp/models"
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccountHandler(t *testing.T) {
	// Connect postgres
	configs.ConnectToDatabase()

	// Initialize the test server
	router := mux.NewRouter()
	router.HandleFunc("/create_account", controllers.CreateAccountHandler)

	// Create a test request
	requestBody := []byte(`{
        "user_id": 1,
        "name": "Test Account",
        "balance": 100.0
    }`)
	req, err := http.NewRequest("POST", "/create_account", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Serve the test request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedResponse := []byte(messages.AccountCreatedSuccessfully)
	if !bytes.Equal(rr.Body.Bytes(), expectedResponse) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedResponse))
	}

	// Check that the account was created in the database
	var account models.Account
	if result := configs.DB.First(&account, "user_id = ?", 1); result.Error != nil {
		t.Errorf("failed to find account in database: %v", result.Error)
	}
	if account.Name != "Test Account" || account.Balance != 100.0 {
		t.Errorf("unexpected account details in database: got %v", account)
	}
}
