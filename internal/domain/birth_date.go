package domain

import (
	"fmt"
	"time"
)

type BirthDate struct {
	value time.Time
}

func NewBirthDate(date string) (*BirthDate, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("data de nascimento inválida: formato deve ser YYYY-MM-DD")
	}

	if parsed.After(time.Now()) {
		return nil, fmt.Errorf("data de nascimento não pode ser no futuro")
	}

	return &BirthDate{value: parsed}, nil
}

func (b BirthDate) String() string {
	return b.value.Format("2006-01-02")
}

func (b BirthDate) Value() time.Time {
	return b.value
}
