package main

import (
	"log"
	"os"

	"go-mymoney/bots"
	"go-mymoney/db"
	"go-mymoney/repository"
	"go-mymoney/service"
	"go-mymoney/utils"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	dsn := os.Getenv("DATABASE_URL")

	if token == "" {
		log.Fatalf("Переменная окружения BOT_TOKEN не задана")
	}
	if dsn == "" {
		log.Fatalf("Переменная окружения DATABASE_URL не задана")
	}

	dsn = utils.EscapePasswordInDSN(dsn)

	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	repo := repository.NewRepository(database)
	svc := service.NewService(repo)

	pref := telebot.Settings{
		Token: token,
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Не удалось создать бота: %v", err)
	}

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Printf("Получено сообщение: %s", c.Text())
		if c.Chat().Type == telebot.ChatPrivate {
			return bots.HandlePrivate(c, svc)
		}
		return bots.HandleText(c, svc)
	})
	bot.Handle(telebot.OnAddedToGroup, func(c telebot.Context) error {
		return bots.HandleAddedToGroup(c, svc)
	})
	bot.Handle(telebot.OnQuery, func(c telebot.Context) error {
		return bots.HandleInlineQuery(c, svc)
	})

	log.Println("Бот запущен")
	bot.Start()
}
