package configs

import (
	"QSuperApp/models"
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	onceDB sync.Once
)

func ConnectToDatabase() {
	onceDB.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_USERNAME"),
				os.Getenv("DATABASE_PASSWORD"),
				os.Getenv("DATABASE_DB"),
			),
		}), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Migrate the schema
		db.AutoMigrate(&models.User{})
		db.AutoMigrate(&models.Order{})
		db.AutoMigrate(&models.Airplane{})
		db.AutoMigrate(&models.Customization{})
		db.AutoMigrate(&models.Authentication{})
		db.AutoMigrate(&models.Payment{})
		db.AutoMigrate(&models.Role{})
		db.AutoMigrate(&models.RoleUser{})
		db.AutoMigrate(&models.Account{})

		// Apply additional migrations
		// if err := db.Exec("ALTER TABLE users ADD COLUMN new_column VARCHAR(255);").Error; err != nil {
		//return err
		///}

		DB = db
	})

}
