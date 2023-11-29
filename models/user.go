package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Cellphone string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	IsAdmin   bool      `gorm:"default:false"`
	// User has many accounts
	Accounts []Account `gorm:"foreignKey:UserID"`
	//Order    []Order   `gorm:"foreignKey:UserID"`

	// User has many RoleUser
	RoleUser []RoleUser `gorm:"foreignKey:UserID"`
}

type GetUserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Cellphone string `json:"cellphone"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Cellphone string `json:"cellphone"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Cellphone string `json:"cellphone"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Cellphone string    `json:"cellphone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsAdmin   bool      `json:"is_admin"`
}

var JWTSecret = []byte("secret")
