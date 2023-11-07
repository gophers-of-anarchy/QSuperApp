package middlewares

import (
	"QSuperApp/messages"
	"QSuperApp/models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for the "/register" and "/login" endpoints
		if r.URL.Path == "/register" || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, messages.Unauthorized, http.StatusUnauthorized)
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return models.JWTSecret, nil
		})
		if err != nil {
			http.Error(w, messages.InvalidToken, http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, messages.InvalidToken, http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			http.Error(w, messages.InvalidToken, http.StatusUnauthorized)
			return
		}
		// Add the authenticated user information to the request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
