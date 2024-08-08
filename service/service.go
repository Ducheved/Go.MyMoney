package service

import (
	"errors"
	"go-mymoney/models"
	"go-mymoney/repository"
	validators "go-mymoney/utils"
	"log"

	"gorm.io/gorm"
)

type Service interface {
	ProcessMessage(userID, chatID int64, message string) error
	GetUserChatsInfo(userID int64) ([]*models.Chat, error)
	GetChatBalance(chatID int64) (float64, error)
	RegisterUserIfNotExists(userID int64) error
	RegisterChatIfNotExists(chatID, userID int64, title string) error
}

type BotService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &BotService{repo: repo}
}

func (s *BotService) ProcessMessage(userID, chatID int64, message string) error {
	amount, currency, err := validators.ValidateMessageFormat(message)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = s.RegisterUserIfNotExists(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currency == "$" || currency == "доллар" || currency == "евро" || currency == "€" {
		amount *= 100
	}

	err = validators.ValidateAmount(amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		chat = &models.Chat{UserID: userID, GroupID: chatID, Balance: 0}
		err = s.repo.SaveChat(chat)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	chat.Balance += amount

	err = s.repo.SaveChat(chat)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *BotService) GetUserChatsInfo(userID int64) ([]*models.Chat, error) {
	return s.repo.GetUserChats(userID)
}

func (s *BotService) GetChatBalance(chatID int64) (float64, error) {
	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		return 0, err
	}
	return chat.Balance, nil
}

func (s *BotService) RegisterUserIfNotExists(userID int64) error {
	tx := s.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user := &models.User{UserID: userID}
			err = s.repo.SaveUser(user)
			if err != nil {
				log.Printf("Ошибка регистрации пользователя: %v", err)
				tx.Rollback()
				return err
			}
		} else {
			log.Printf("Ошибка получения пользователя: %v", err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (s *BotService) RegisterChatIfNotExists(chatID, userID int64, title string) error {
	tx := s.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := s.RegisterUserIfNotExists(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = s.repo.GetChatByID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			chat := &models.Chat{UserID: userID, GroupID: chatID, Balance: 0}
			err = s.repo.RegisterChat(chat)
			if err != nil {
				log.Printf("Ошибка регистрации чата: %v", err)
				tx.Rollback()
				return err
			}
		} else {
			log.Printf("Ошибка получения чата: %v", err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
