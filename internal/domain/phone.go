package domain

import (
	"errors"
	"regexp"
)

type Phone struct {
	Number string
}

func NewPhone(number string) (*Phone, error) {
	if !isValidPhone(number) {
		return nil, errors.New("invalid phone number format")
	}
	return &Phone{Number: number}, nil
}

func isValidPhone(number string) bool {
	// Example regex for phone validation (adjust as needed)
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(number)
}
