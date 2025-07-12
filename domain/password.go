package domain

import (
	"fmt"
	"unicode"
)

type Password string

func NewPassword(value string) (*Password, error) {
	if err := validatePassword(value); err != nil {
		return nil, err
	}

	p := Password(value)
	return &p, nil
}

func validatePassword(s string) error {
	if len(s) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return fmt.Errorf("password must contain at least one uppercase letter")
	case !hasLower:
		return fmt.Errorf("password must contain at least one lowercase letter")
	case !hasDigit:
		return fmt.Errorf("password must contain at least one digit")
	case !hasSpecial:
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func (p *Password) Equals(value string) bool {
	return *p == Password(value)
}
