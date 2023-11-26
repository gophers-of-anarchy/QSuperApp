package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/functions"
	"QSuperApp/messages"
	"QSuperApp/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(ctx echo.Context) error {
	req := models.RegisterRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.FailedPasswordHashGeneration)
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
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateUser)
	}

	accessToken, refreshToken, err := functions.GenerateTokens(user.Username, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAuthTokens)
	}

	response := models.RegisterResponse{AccessToken: accessToken, RefreshToken: refreshToken}

	return ctx.JSON(
		http.StatusOK,
		response,
	)
}

func LoginHandler(ctx echo.Context) error {
	req := models.LoginRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	var user models.User
	result := configs.DB.Where("username = ?", req.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusUnauthorized, messages.UsernameOrPasswordIncorrect)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.UsernameOrPasswordIncorrect)
	}
	accessToken, refreshToken, err := functions.GenerateTokens(user.Username, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAuthTokens)
	}
	response := models.RegisterResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	return ctx.JSON(
		http.StatusOK,
		response,
	)
}

func RestrictedAreaHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)
	msg := fmt.Sprintf(messages.Welcome, userID)
	w.Write([]byte(msg))
}
