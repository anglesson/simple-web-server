package client_domain

import (
	"errors"

	common_entities "github.com/anglesson/simple-web-server/internal/common/domain/entities"
	common_vo "github.com/anglesson/simple-web-server/internal/common/domain/valueobjects"
)

type Client struct {
	ID       string
	Name     string
	Cpf      *common_vo.CPF
	Birthday *common_vo.BirthDay
	contact  *common_entities.Contact
}

func NewClient(id, name, cpf, birthday, email, phone string) (*Client, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	voCpf, err := common_vo.NewCPF(cpf)
	if err != nil {
		return nil, err
	}

	voBirthDay, err := common_vo.NewBirthDay(2020, 9, 2)
	if err != nil {
		return nil, err
	}

	contactEntity, err := common_entities.NewContact(email, phone)
	if err != nil {
		return nil, err
	}

	return &Client{
		ID:       id,
		Name:     name,
		Cpf:      voCpf,
		Birthday: voBirthDay,
		contact:  contactEntity,
	}, nil
}

func (c *Client) UpdateContact(email, phone string) error {
	contactEntity, err := common_entities.NewContact(email, phone)
	if err != nil {
		return err
	}
	c.contact = contactEntity
	return nil
}
