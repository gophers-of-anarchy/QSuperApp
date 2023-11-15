package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(messages.RegisteredSuccessfully))
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

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(models.JWTSecret)
	if err != nil {
		http.Error(w, messages.FailedToCreateToken, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func RestrictedAreaHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)
	msg := fmt.Sprintf(messages.Welcome, userID)
	w.Write([]byte(msg))
}
