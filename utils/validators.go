package utils

import (
	"errors"
	"regexp"
	"strconv"
)

var validCurrencies = map[string]bool{
	"руб":      true,
	"рубль":    true,
	"рублей":   true,
	"рубли":    true,
	"$":        true,
	"доллар":   true,
	"долларов": true,
	"€":        true,
	"евро":     true,
}

func ValidateMessageFormat(message string) (float64, string, error) {
	re := regexp.MustCompile(`([+-]?\d+(\.\d+)?)([a-zA-Zа-яА-Я$€]+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) != 4 {
		return 0, "", errors.New("invalid message format")
	}

	if ValidateURL(message) {
		return 0, "", errors.New("message contains URL")
	}

	// if ValidateSQLInjection(message) {
	// 	return 0, "", errors.New("message contains SQL injection attempt")
	// }

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", err
	}

	currency := matches[3]
	if !validCurrencies[currency] {
		return 0, "", errors.New("invalid currency")
	}

	if err := ValidateAmount(amount); err != nil {
		return 0, "", err
	}

	return amount, currency, nil
}

func ValidateAmount(amount float64) error {
	if amount > 10000 || amount < -10000 {
		return errors.New("amount exceeds the maximum limit of ±10000")
	}
	return nil
}

func ValidateSQLInjection(text string) bool {
	// Исключаем допустимые символы валютных операций
	sqlRe := regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|EXEC|UNION|--|;|')`)
	return sqlRe.MatchString(text) && !regexp.MustCompile(`^[+-]?\d+(\.\d+)?[a-zA-Zа-яА-Я$€]+$`).MatchString(text)
}

func ValidateURL(text string) bool {
	urlRe := regexp.MustCompile(`https?://[^\s]+`)
	return urlRe.MatchString(text)
}
