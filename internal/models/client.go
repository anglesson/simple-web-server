package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name           string     `json:"name"`
	CPF            string     `gorm:"unique" json:"cpf"`
	DataNascimento string     `json:"data_nascimento"`
	ContactID      uint       `json:"contact_id"`
	Validated      bool       `json:"validated"`
	Contact        Contact    `gorm:"foreignKey:ContactID"`
	Creators       []*Creator `gorm:"many2many:client_creators"`
}

func NewClient(name, cpf, dataNascimento, email, phone string, creator *Creator) *Client {
	contact := NewContact(email, phone)
	return &Client{
		Name:           name,
		CPF:            cpf,
		DataNascimento: dataNascimento,
		Contact:        contact,
		Validated:      false,
		Creators: []*Creator{
			creator,
		},
	}
}

func (c *Client) Update(name, cpf, email, phone string) {
	c.Name = name
	c.CPF = cpf
	c.Contact.Email = email
	c.Contact.Phone = phone
}
