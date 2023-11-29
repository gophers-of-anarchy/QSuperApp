package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateAccountHandler(ctx echo.Context) error {

	// I think when we check the authorization using middleware
	// we dont need the 'ok' and we ignore it
	userID, _ := ctx.Get("user_id").(uint)

	req := models.AccountCreateRequest{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	account := models.Account{
		UserID:  userID,
		Name:    req.Name,
		Balance: req.Balance,
	}

	if result := configs.DB.Create(&account); result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAccount)
	}

	response := models.CreateAccountResponse{
		Message: messages.AccountCreatedSuccessfully,
		Data:    models.AccountResponseData{Name: req.Name, Balance: fmt.Sprintf("%v", req.Balance)},
	}

	return ctx.JSON(
		http.StatusCreated,
		response,
	)

}

func UpdateAccount(ctx echo.Context) error {

	userID, _ := ctx.Get("user_id").(uint)

	req := models.UpdateAccountRequest{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	accountID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}

	var account models.Account

	if err := configs.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.AccountNotFound)
	}

	account.Name = req.Name

	if err := configs.DB.Save(&account).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.UpdateAccountFailed)
	}

	response := models.CreateAccountResponse{
		Message: messages.AccountUpdatedSuccessfully,
		Data:    models.AccountResponseData{Name: req.Name, Balance: fmt.Sprintf("%v", account.Balance)},
	}

	return ctx.JSON(
		http.StatusOK,
		response,
	)
}
