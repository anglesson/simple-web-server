package domain

import (
	"errors"
	"time"
)

type Creator struct {
	Name      string
	CPF       *CPF
	Birthdate *BirthDay
	Contact   *Contact
}

func NewCreator(name, email, cpf, phone, birthdate string) (*Creator, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}

	if len(name) > 255 {
		name = name[0:255]
	}

	contactVo, err := NewContact(email, phone)
	if err != nil {
		return nil, err
	}

	cpfVo, err := NewCPF(cpf)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return nil, err
	}

	year := parsedDate.Year()
	month := int(parsedDate.Month())
	day := parsedDate.Day()

	birthdateVo, err := NewBirthDay(year, month, day)
	if err != nil {
		return nil, err
	}

	return &Creator{
		Name:      name,
		CPF:       cpfVo,
		Birthdate: birthdateVo,
		Contact:   contactVo,
	}, nil
}
