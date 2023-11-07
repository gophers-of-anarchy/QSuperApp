package configs

import (
	"QSuperApp/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
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

		DB = db
	})

}
