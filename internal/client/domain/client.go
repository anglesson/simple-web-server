package client_domain

import (
	"errors"
	"fmt"

	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
)

const (
	MinNameLength = 5
	MaxNameLength = 255
)

type Client struct {
	Name     string
	Phone    string
	BirthDay *common_domain.BirthDate
	CPF      common_domain.CPF
	Email    string
}

func NewClient(name, cpf, birthDay, email, phone string) (*Client, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(name) < 5 || len(name) > 255 {
		return nil, fmt.Errorf("o nome dever ter entre %v e %v caracteres", MinNameLength, MaxNameLength)
	}

	validCPF, err := common_domain.NewCPF(cpf)
	if err != nil {
		return nil, fmt.Errorf("CPF inválido para o cliente: %w", err)
	}

	birth, err := common_domain.NewBirthDate(birthDay)
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
