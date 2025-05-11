package models

import "gorm.io/gorm"

type Client struct {
	*gorm.Model
	Name      string     `json:"name"`
	CPF       string     `gorm:"unique" json:"cpf"`
	ContactID uint       `json:"contact_id"`
	Contact   Contact    `gorm:"foreignKey:ContactID"`
	Creators  []*Creator `gorm:"many2many:client_creator"`
}

func NewClient(name, cpf, email, phone string, creator *Creator) *Client {
	contact := NewContact(email, phone)
	return &Client{
		Name:    name,
		CPF:     cpf,
		Contact: contact,
		Creators: []*Creator{
			creator,
		},
	}
}
