package domain

import (
	"errors"
)

type Client struct {
	ID       string
	Name     string
	Cpf      *CPF
	Birthday *BirthDay
	contact  *Contact
}

func NewClient(id, name, cpf, birthday, email, phone string) (*Client, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	voCpf, err := NewCPF(cpf)
	if err != nil {
		return nil, err
	}

	voBirthDay, err := NewBirthDay(2020, 9, 2)
	if err != nil {
		return nil, err
	}

	contactEntity, err := NewContact(email, phone)
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
	contactEntity, err := NewContact(email, phone)
	if err != nil {
		return err
	}
	c.contact = contactEntity
	return nil
}
