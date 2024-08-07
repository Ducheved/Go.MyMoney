package db

import (
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-mymoney/models"
)

func NewPostgresDB(dbURL string) (*gorm.DB, error) {
	if !strings.Contains(dbURL, "search_path") {
		dbURL += "?search_path=public"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Chat{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
