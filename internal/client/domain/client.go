package domain

import "errors"

type Client struct {
	Name     string
	Phone    string
	BirthDay string
	CPF      string
	Email    string
}

func NewClient(name, cpf, birthDay, email, phone string) (*Client, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if cpf == "" {
		return nil, errors.New("CPF is required")
	}
	if birthDay == "" {
		return nil, errors.New("birthday is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if phone == "" {
		return nil, errors.New("phone is required")
	}
	return &Client{
		Name:     name,
		CPF:      cpf,
		Email:    email,
		Phone:    phone,
		BirthDay: birthDay,
	}, nil
}
