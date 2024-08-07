package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID int64 `gorm:"unique_index"`
}

type Chat struct {
	gorm.Model
	GroupID int64 `gorm:"unique_index"`
	UserID  int64
	Balance float64
}
