package domain

import (
	"errors"
	"fmt"
)

const (
	MinNameLength = 5
	MaxNameLength = 255
)

type Client struct {
	ID       uint
	Name     string
	Phone    string
	BirthDay *BirthDate
	CPF      CPF
	Email    string
}

func NewClient(name, cpf, birthDay, email, phone string) (*Client, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(name) < 5 || len(name) > 255 {
		return nil, fmt.Errorf("o nome dever ter entre %v e %v caracteres", MinNameLength, MaxNameLength)
	}

	validCPF, err := NewCPF(cpf)
	if err != nil {
		return nil, fmt.Errorf("CPF inválido para o cliente: %w", err)
	}

	birth, err := NewBirthDate(birthDay)
	if err != nil {
		return nil, err
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if phone == "" {
		return nil, errors.New("phone is required")
	}
	return &Client{
		Name:     name,
		CPF:      validCPF,
		Email:    email,
		Phone:    phone,
		BirthDay: birth,
	}, nil
}

func (c *Client) Update(name, cpf, email, phone string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) < 5 || len(name) > 255 {
		return fmt.Errorf("o nome dever ter entre %v e %v caracteres", MinNameLength, MaxNameLength)
	}

	validCPF, err := NewCPF(cpf)
	if err != nil {
		return fmt.Errorf("CPF inválido para o cliente: %w", err)
	}

	if email == "" {
		return errors.New("email is required")
	}
	if phone == "" {
		return errors.New("phone is required")
	}

	c.Name = name
	c.CPF = validCPF
	c.Email = email
	c.Phone = phone

	return nil
}
