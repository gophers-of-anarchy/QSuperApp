package middlewares

import (
	"QSuperApp/messages"
	"QSuperApp/models"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/register" || c.Path() == "/login" {
			return next(c)
		}
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.Unauthorized})
		}
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return models.JWTSecret, nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.InvalidToken})
		}
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.InvalidToken})
		}
		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": messages.InvalidToken})
		}

		// Add the authenticated user information to the Echo context
		c.Set("user_id", claims.UserID)

		return next(c)

	}

}
