package repository

import (
	"go-mymoney/models"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	SaveChat(chat *models.Chat) error
	GetChatByID(chatID int64) (*models.Chat, error)
	UpdateBalance(chatID int64, amount float64) error
	GetUserChats(userID int64) ([]*models.Chat, error)
	SaveUser(user *models.User) error
	GetUserByID(userID int64) (*models.User, error)
	RegisterChat(chat *models.Chat) error
	BeginTransaction() *gorm.DB
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) SaveChat(chat *models.Chat) error {
	err := r.db.Save(chat).Error
	if err != nil {
		log.Printf("Ошибка сохранения чата: %v", err)
	}
	return err
}

func (r *GormRepository) GetChatByID(chatID int64) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.First(&chat, "group_id = ?", chatID).Error
	if err != nil {
		log.Printf("Ошибка получения чата по ID: %v", err)
	}
	return &chat, err
}

func (r *GormRepository) UpdateBalance(chatID int64, amount float64) error {
	return r.db.Model(&models.Chat{}).Where("group_id = ?", chatID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *GormRepository) GetUserChats(userID int64) ([]*models.Chat, error) {
	var chats []*models.Chat
	err := r.db.Where("user_id = ?", userID).Find(&chats).Error
	return chats, err
}

func (r *GormRepository) SaveUser(user *models.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		log.Printf("Ошибка сохранения пользователя: %v", err)
	}
	return err
}

func (r *GormRepository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "user_id = ?", userID).Error
	if err != nil {
		log.Printf("Ошибка получения пользователя по ID: %v", err)
	}
	return &user, err
}

func (r *GormRepository) RegisterChat(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

func (r *GormRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
