package middlewares

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsPassengerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := false

		userID, ok := c.Get("user_id").(uint)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.Unauthorized})
		}

		var user = models.User{ID: userID}
		result := configs.DB.Preload("RoleUser.Role").First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": messages.UserNotFound})
		}

		for _, roleUser := range user.RoleUser {
			if roleUser.Role.RoleTitle == "PS" {
				access = true
			}
		}

		if !access {
			return c.JSON(http.StatusForbidden, map[string]string{"error": messages.ErrUnauthorizedAccess})
		}
		return next(c)

	}
}

func AuthAndIsPassengerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthMiddleware(IsPassengerMiddleware(next))
}
