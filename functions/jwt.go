package functions

import (
	"QSuperApp/models"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokens(username string, userID uint) (string, string, error) {

	// ACCESS
	accessClaims := models.Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// this should be store in .env I should ask arman later
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSignedToken, err := accessToken.SignedString(models.JWTSecret)
	if err != nil {
		log.Println("(GenerateAccessToken) Error :", err)
		return "", "", err
	}

	// REFRESH
	refreshClaims := models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSignedToken, err := refreshToken.SignedString(models.JWTSecret)
	if err != nil {
		log.Println("(GenerateRefreshToken) Error :", err)
		return "", "", err
	}

	return accessSignedToken, refreshSignedToken, nil
}
