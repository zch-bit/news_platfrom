package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"newsplatform/internal/models"
)

var DB *gorm.DB

func ConnectDB() error {
	db, err := gorm.Open(sqlite.Open("articles.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database")
	}

	err = db.AutoMigrate(&models.News{})
	if err != nil {
		return fmt.Errorf("failed to migrate database")
	}
	DB = db
	return nil
}
