package bots

import (
	"fmt"
	"go-mymoney/service"
	validators "go-mymoney/utils"
	"log"
	"regexp"
	"strconv"

	"gopkg.in/telebot.v3"
)

func HandleText(c telebot.Context, svc service.Service) error {
	userID := c.Sender().ID
	chatID := c.Chat().ID

	if c.Text() == "/balance" {
		balance, err := svc.GetChatBalance(chatID)
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка получения баланса: %v", err))
		}
		return c.Send(fmt.Sprintf("Текущий баланс: %.2f", balance))
	}

	if len(c.Text()) >= len("/gen ") && c.Text()[:len("/gen ")] == "/gen " {
		return handleGenCommand(c)
	}

	amount, currency, err := validators.ValidateMessageFormat(c.Text())
	if err != nil {
		// return c.Send(fmt.Sprintf("Ошибка валидации сообщения: %v", err))
		return nil
	}

	sign := "+"
	if amount < 0 {
		sign = "-"
	}

	action := "добавили"
	if sign == "-" {
		action = "вычли"
	}

	err = svc.ProcessMessage(userID, chatID, fmt.Sprintf("%s%.2f%s", sign, amount, currency))
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка обработки сообщения: %v", err))
	}

	balance, err := svc.GetChatBalance(chatID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка получения баланса: %v", err))
	}

	log.Printf("Вы %s: %.2f %s для чата %d", action, amount, currency, chatID)
	return c.Send(fmt.Sprintf("Вы %s: %.2f %s. Текущий баланс: %.2f %s", action, amount, currency, balance, currency))
}

func HandleAddedToGroup(c telebot.Context, svc service.Service) error {
	chatID := c.Chat().ID
	userID := c.Sender().ID

	err := svc.RegisterChatIfNotExists(chatID, userID, c.Chat().Title)
	if err != nil {
		log.Printf("Ошибка регистрации чата: %v", err)
		return err
	}

	log.Printf("Чат зарегистрирован: %d", chatID)
	return c.Send("Добро пожаловать, я ваш банк!")
}

func HandlePrivate(c telebot.Context, svc service.Service) error {
	userID := c.Sender().ID

	switch c.Text() {
	case "/balance":
		chats, err := svc.GetUserChatsInfo(userID)
		if err != nil {
			log.Printf("Ошибка получения информации о чатах: %v", err)
			return err
		}

		if len(chats) == 0 {
			return c.Send("У вас нет зарегистрированных чатов.")
		}

		var response string
		for _, chat := range chats {
			response += fmt.Sprintf("Чат ID: %d, Баланс: %.2f₽\n", chat.GroupID, chat.Balance)
		}

		log.Printf("Отправка информации о чатах пользователю %d", userID)
		return c.Send(response)
	case "/mygroups":
		chats, err := svc.GetUserChatsInfo(userID)
		if err != nil {
			log.Printf("Ошибка получения информации о чатах: %v", err)
			return err
		}

		if len(chats) == 0 {
			return c.Send("У вас нет зарегистрированных чатов.")
		}

		var response string
		for _, chat := range chats {
			response += fmt.Sprintf("Чат ID: %d, Баланс: %.2f₽\n", chat.GroupID, chat.Balance)
		}

		log.Printf("Отправка информации о чатах пользователю %d", userID)
		return c.Send(response)
	default:
		return c.Send("Неизвестная команда")
	}
}

func HandleInlineQuery(c telebot.Context, svc service.Service) error {
	query := c.Query().Text
	userID := c.Sender().ID
	chatID := c.Chat().ID

	re := regexp.MustCompile(`([+-])(\d+)([a-zA-Zа-яА-Я$€]*)`)
	matches := re.FindStringSubmatch(query)
	log.Printf("Обработка запроса: %s, найденные совпадения: %v", query, matches)
	if len(matches) == 4 {
		sign := matches[1]
		amount, _ := strconv.ParseFloat(matches[2], 64)
		if sign == "-" {
			amount = -amount
		}
		err := svc.ProcessMessage(userID, chatID, sign+matches[2]+matches[3])
		if err != nil {
			log.Printf("Ошибка обработки сообщения: %v", err)
			return err
		}
		currentBalance, err := svc.GetChatBalance(chatID)
		if err != nil {
			log.Printf("Ошибка получения баланса: %v", err)
			return err
		}
		log.Printf("Баланс обновлен: %.2f для чата %d", amount, chatID)
		results := []telebot.Result{
			&telebot.ArticleResult{
				Title:       "Баланс обновлен",
				Description: fmt.Sprintf("Баланс обновлен: %.2f₽. Текущий баланс: %.2f₽", amount, currentBalance),
				Text:        fmt.Sprintf("Баланс обновлен: %.2f₽. Текущий баланс: %.2f₽", amount, currentBalance),
			},
		}
		return c.Answer(&telebot.QueryResponse{
			Results: results,
		})
	} else {
		log.Printf("Сообщение не распознано: %s", query)
	}

	return nil
}
