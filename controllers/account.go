package controllers

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateAccountHandler(ctx echo.Context) error {

	// I think when we check the authorization using middleware
	// we dont need the 'ok' and we ignore it
	userID, _ := ctx.Get("user_id").(uint)

	req := models.AccountCreateRequest{}

	var checkUser = models.Account{UserID: userID}
	foundAccount := configs.DB.First(&checkUser)
	if foundAccount == nil {
		return ctx.JSON(http.StatusBadRequest, messages.UserAlreadyHasAccount)
	}

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
		http.StatusOK,
		response,
	)

}
