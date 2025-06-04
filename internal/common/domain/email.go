package common_domain

import (
	"errors"
	"regexp"
)

type Email struct {
	value string
}

func NewEmail(address string) (Email, error) {
	// Validate email format using a regular expression.
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(address) {
		return Email{}, errors.New("endereço de email inválido")
	}
	return Email{value: address}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

func (e Email) String() string {
	return e.value
}
