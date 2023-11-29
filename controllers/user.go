package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/functions"
	"QSuperApp/messages"
	"QSuperApp/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

func UpdateUserHandler(ctx echo.Context) error {

	userID, _ := ctx.Get("user_id").(uint)

	req := models.UpdateUserRequest{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	var user models.User

	if err := configs.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.AccountNotFound)
	}

	user.Cellphone = req.Cellphone
	user.Email = req.Email
	user.Username = req.Username

	if err := configs.DB.Save(&user).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.UpdateAccountFailed)
	}

	response := models.UpdateUserResponse{
		Message: messages.UserUpdatedSuccefully,
	}

	return ctx.JSON(
		http.StatusOK,
		response,
	)
}

func GetUserHandler(ctx echo.Context) error {
	userID, _ := ctx.Get("user_id").(uint)

	var user = models.User{ID: userID}
	result := configs.DB.First(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusBadRequest, messages.UserNotFound)
	}

	response := models.GetUserResponse{
		Username:  user.Username,
		Email:     user.Email,
		Cellphone: user.Cellphone,
	}
	return ctx.JSON(
		http.StatusOK,
		response,
	)

}

func GetUserByIdHandler(ctx echo.Context) error {
	userId, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var user models.User
	if err := configs.DB.Preload("Accounts").Preload("RoleUser").First(&user, userId).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, messages.UserNotFound)
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Cellphone: user.Cellphone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsAdmin:   user.IsAdmin,
	}

	return ctx.JSON(http.StatusOK, response)

}
