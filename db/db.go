package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-mymoney/models"
)

func NewPostgresDB(dbURL string) (*gorm.DB, error) {
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
