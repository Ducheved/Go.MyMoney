package validators

import (
	"errors"
	"regexp"
	"strconv"
)

func ValidateMessageFormat(message string) (float64, string, error) {
	re := regexp.MustCompile(`([+-]\d+)([a-zA-Zа-яА-Я$€]+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) != 3 {
		return 0, "", errors.New("invalid message format")
	}

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", err
	}

	if matches[2] == "" {
		return 0, "", errors.New("currency must be specified")
	}

	return amount, matches[2], nil
}

func ValidateAmount(amount float64) error {
	if amount > 10000 || amount < -10000 {
		return errors.New("amount exceeds the maximum limit of ±10000")
	}
	return nil
}

func ValidateSQLInjection(text string) bool {
	sqlRe := regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|EXEC|UNION|--|;|')`)
	return sqlRe.MatchString(text)
}

func ValidateURL(text string) bool {
	urlRe := regexp.MustCompile(`https?://[^\s]+`)
	return urlRe.MatchString(text)
}
