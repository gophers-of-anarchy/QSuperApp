package middlewares

import (
	"QSuperApp/configs"
	"QSuperApp/messages"
	"QSuperApp/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, ok := c.Get("user_id").(uint)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.Unauthorized})
		}

		var user = models.User{ID: userID}
		result := configs.DB.First(&user)

		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		if !user.IsAdmin {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": messages.InvalidUserType})
		}

		return next(c)
	}
}

func AuthAndAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthMiddleware(IsAdminMiddleware(next))
}
