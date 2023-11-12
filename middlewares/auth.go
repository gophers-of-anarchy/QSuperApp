package middlewares

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Message string
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	var response Response
	return func(ctx echo.Context) error {
		if ctx.Request().URL.Path == "/Register" || ctx.Request().URL.Path == "/Login" || ctx.Request().URL.Path == "/Email_Verification" {
			return next(ctx)
		}
		tokenString := ctx.Request().Header.Get("Authorization")
		if tokenString == "" {
			response.Message = messages.Unauthorized
			ctx.JSON(http.StatusUnauthorized, response)
			return errors.New(strings.ToLower(messages.Unauthorized))
		}
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SecretKey")), nil
		})
		if err != nil {
			response.Message = messages.InvalidToken
			ctx.JSON(http.StatusUnauthorized, response)
			return errors.New(strings.ToLower(messages.InvalidToken))
		}
		if !token.Valid {
			response.Message = messages.InvalidToken
			ctx.JSON(http.StatusUnauthorized, response)
			return errors.New(strings.ToLower(messages.InvalidToken))
		}
		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			response.Message = messages.InvalidToken
			ctx.JSON(http.StatusUnauthorized, response)
			return errors.New(strings.ToLower(messages.InvalidToken))
		}
		if strings.Contains(ctx.Request().URL.Path, "/admin") {
			var users []models.User
			result := configs.DB.Find(&users)
			if result.Error != nil {
				response.Message = messages.InternalError
				ctx.JSON(http.StatusInternalServerError, response)
				return result.Error
			}
			for _, user := range users {
				if user.ID == claims.UserID {
					if user.Role != "admin" {
						response.Message = messages.Unauthorized
						ctx.JSON(http.StatusUnauthorized, response)
						return errors.New(strings.ToLower(messages.Unauthorized))
					}
					ctx.Set("user_id", claims.UserID)
					response.Message = messages.AuthenticatedSuccessfully
					ctx.JSON(http.StatusOK, response)
					return next(ctx)
				}
			}
		}
		ctx.Set("user_id", claims.UserID)
		response.Message = messages.AuthenticatedSuccessfully
		ctx.JSON(http.StatusOK, response)
		return next(ctx)
	}
}
