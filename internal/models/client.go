package models

import "gorm.io/gorm"

type Client struct {
	*gorm.Model
	Name      string  `json:"name"`
	CPF       string  `json:"cpf"`
	ContactID uint    `json:"contact_id"`
	Contact   Contact `gorm:"foreignKey:ContactID"`
}

func NewClient(name, cpf, email, phone string) *Client {
	contact := NewContact(email, phone)
	return &Client{
		Name:    name,
		CPF:     cpf,
		Contact: contact,
	}
}
