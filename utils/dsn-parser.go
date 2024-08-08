package utils

import (
	"log"
	"net/url"
	"strings"
)

func EscapePasswordInDSN(dsn string) string {
	start := strings.Index(dsn, "://") + 3
	end := strings.Index(dsn[start:], "@")
	if end == -1 {
		log.Fatalf("Не удалось найти пароль в строке DSN")
	}
	end += start

	userInfo := dsn[start:end]
	rest := dsn[end:]

	userPass := strings.SplitN(userInfo, ":", 2)
	if len(userPass) != 2 {
		log.Fatalf("Не удалось разделить имя пользователя и пароль в строке DSN")
	}
	username := userPass[0]
	password := userPass[1]

	escapedPassword := url.QueryEscape(password)

	escapedUserInfo := username + ":" + escapedPassword
	dsn = dsn[:start] + escapedUserInfo + rest

	_, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("Не удалось распарсить DSN: %v", err)
	}

	return dsn
}
