package service

import (
	"errors"
	"log"
	"regexp"
	"strconv"

	"go-mymoney/models"
	"go-mymoney/repository"

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
	err := s.RegisterUserIfNotExists(userID)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`([+-]\d+)([a-zA-Zа-яА-Я$€]+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) != 3 {
		return errors.New("invalid message format")
	}

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return err
	}

	if matches[2] == "" {
		return errors.New("currency must be specified")
	}

	if matches[2] == "$" || matches[2] == "доллар" || matches[2] == "евро" || matches[2] == "€" {
		amount *= 100
	}

	if amount > 10000 || amount < -10000 {
		return errors.New("amount exceeds the maximum limit of ±10000")
	}

	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		chat = &models.Chat{UserID: userID, GroupID: chatID, Balance: 0}
		err = s.repo.SaveChat(chat)
		if err != nil {
			return err
		}
	}

	chat.Balance += amount

	err = s.repo.SaveChat(chat)
	if err != nil {
		return err
	}

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
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user := &models.User{UserID: userID}
			err = s.repo.SaveUser(user)
			if err != nil {
				log.Printf("Ошибка регистрации пользователя: %v", err)
				return err
			}
		} else {
			log.Printf("Ошибка получения пользователя: %v", err)
			return err
		}
	}
	return nil
}

func (s *BotService) RegisterChatIfNotExists(chatID, userID int64, title string) error {
	err := s.RegisterUserIfNotExists(userID)
	if err != nil {
		return err
	}

	_, err = s.repo.GetChatByID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			chat := &models.Chat{UserID: userID, GroupID: chatID, Balance: 0}
			err = s.repo.RegisterChat(chat)
			if err != nil {
				log.Printf("Ошибка регистрации чата: %v", err)
				return err
			}
		} else {
			log.Printf("Ошибка получения чата: %v", err)
			return err
		}
	}
	return nil
}
