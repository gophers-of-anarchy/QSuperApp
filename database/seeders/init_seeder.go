package main

import (
	"QSuperApp/configs"
	"QSuperApp/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	configs.ConnectToDatabase()

	// Migrate the User model to create the "users" table
	configs.DB.AutoMigrate(&models.User{})
	//configs.DB.AutoMigrate(&models.Order{})

	// Seed the database with 3 users
	for i := 0; i < 3; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := models.User{
			Username:  fmt.Sprintf("user%d", i),
			Password:  string(hashedPassword),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Cellphone: fmt.Sprintf("123-456-789%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		configs.DB.Create(&user)

		// Seed the database with accounts for each user
		account := models.Account{
			UserID:    user.ID,
			Name:      "Savings",
			Balance:   1000.0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		configs.DB.Create(&account)

		// Seed the database with cards for each account
		card := models.Card{
			AccountID: account.ID,
			Number:    fmt.Sprintf("123456789012%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		configs.DB.Create(&card)
	}
}
