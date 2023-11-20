package config

import (
	"fmt"
	"github.com/Hafidzurr/project3_group2_glng-ks-08/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres dbname=Kanban_Board_db password=Hafidzurr1 sslmode=disable" // Lokal DB config
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database. Error:", err)
		return nil, err
	}

	// Call Migrate function to perform database migrations
	if err := Migrate(db); err != nil {
		fmt.Println("Failed to perform database migrations. Error:", err)
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Category{}, &models.Task{})
}
