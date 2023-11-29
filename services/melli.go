package services

import (
	"QSuperApp/helpers"
	"QSuperApp/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type MelliPaymentService struct {
	TerminalKey string
}

func (s *MelliPaymentService) SETTermianlKey(apikey string) {
	s.TerminalKey = apikey
}

func Payment(Amount float64, orderID uint) (models.MelliPaymentResponse, error) {
	var responseObj models.MelliPaymentResponse

	iranLocation, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Println("Failed to load Iran location:", err)
		return responseObj, err
	}

	// Get the current time in Iran timezone
	iranTime := time.Now().In(iranLocation)
	formattedTime := iranTime.Format("01/02/2006 3:04:05 PM")

	requestData := models.PaymentRequest{
		MerchantId:    "1186",
		TerminalId:    "7NDIXXW0",
		ReturnUrl:     "http://localhost:8080/api/v1/payment/verifypage",
		LocalDateTime: formattedTime,
		OrderId:       orderID,
		Amount:        int(Amount),
	}

	// Create the data string for encryption
	data := fmt.Sprintf("%s;%d;%d", requestData.TerminalId, requestData.OrderId, requestData.Amount)

	key := os.Getenv("TERMINAL_KEY")
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Println("Error decoding key:", err)
		return responseObj, err
	}

	// Create the cipher block
	encrypted, err := helpers.TripleEcbDesEncrypt([]byte(data), decodedKey)
	if err != nil {
		log.Println("Error encrypting data:", err)
		return responseObj, err
	}

	requestData.SignData = encrypted

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Failed to marshal request body:", err)
		return responseObj, err
	}

	req, err := http.NewRequest("POST", os.Getenv("PAYMENT_REQUEST_URL"), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to create request:", err)
		return responseObj, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return responseObj, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return responseObj, err
	}

	err = json.Unmarshal(respBody, &responseObj)
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		return responseObj, err
	}
	responseObj.URL = fmt.Sprintf("%s?token=%s", os.Getenv("PAYMENT_GATEWAY_URL"), responseObj.Token)
	return responseObj, nil
}

func VerifyPayment(token string) (models.MelliVerifyPaymentResponse, error) {
	var responseObj models.MelliVerifyPaymentResponse
	key := os.Getenv("TERMINAL_KEY")
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Println("Error decoding key:", err)
		return responseObj, err
	}

	// Create the cipher block
	encrypted, err := helpers.TripleEcbDesEncrypt([]byte(token), decodedKey)
	if err != nil {
		log.Println("Error encrypting data:", err)
		return responseObj, err
	}
	requestData := models.VerifyPaymentRequest{
		SignData: encrypted,
		Token:    token,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Failed to marshal request body:", err)
		return responseObj, err
	}

	req, err := http.NewRequest("POST", os.Getenv("PAYMENT_VERIFY_URL"), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to create request:", err)
		return responseObj, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return responseObj, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return responseObj, err
	}

	err = json.Unmarshal(respBody, &responseObj)
	fmt.Println(string(responseObj.Description))
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		return responseObj, err
	}
	return responseObj, nil
}
