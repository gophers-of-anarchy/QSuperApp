package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/functions"
	"QSuperApp/messages"
	"QSuperApp/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req models.RegisterRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, messages.FailedPasswordHashGeneration, http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		Cellphone: req.Cellphone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = configs.DB.Create(&user).Error
	if err != nil {
		http.Error(w, messages.FailedToCreateUser, http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := functions.GenerateTokens(user.Username, user.ID)
	if err != nil {
		http.Error(w, messages.FailedToCreateAuthTokens, http.StatusInternalServerError)
	}

	response := models.RegisterResponse{AccessToken: accessToken, RefreshToken: refreshToken}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(messages.RegisteredSuccessfully))
	json.NewEncoder(w).Encode(response)

	log.Printf("user with username: %v created at %v \n", user.Username, time.Now())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req models.LoginRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, messages.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	var user models.User
	result := configs.DB.Where("username = ?", req.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, messages.UsernameOrPasswordIncorrect, http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, messages.UsernameOrPasswordIncorrect, http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := functions.GenerateTokens(user.Username, user.ID)
	if err != nil {
		http.Error(w, messages.FailedToCreateAuthTokens, http.StatusInternalServerError)
	}

	response := models.RegisterResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	json.NewEncoder(w).Encode(response)
}

func RestrictedAreaHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)
	msg := fmt.Sprintf(messages.Welcome, userID)
	w.Write([]byte(msg))
}
