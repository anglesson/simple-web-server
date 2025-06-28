package models

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name      string     `json:"name"`
	CPF       string     `gorm:"unique" json:"cpf"`
	Birthdate string     `json:"birthdate"`
	ContactID uint       `json:"contact_id"`
	Validated bool       `json:"validated"`
	Contact   Contact    `gorm:"foreignKey:ContactID"`
	Creators  []*Creator `gorm:"many2many:client_creators"`
	Purchases []*Purchase
}

func NewClient(name, cpf, birthDate, email, phone string, creator *Creator) *Client {
	contact := NewContact(email, phone)
	return &Client{
		Name:      name,
		CPF:       cpf,
		Birthdate: birthDate,
		Contact:   contact,
		Validated: false,
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

func (c *Client) TotalPurchasesByEbook(ebookID uint) int {
	var count int
	for _, purchase := range c.Purchases {
		if purchase.EbookID == ebookID {
			count++
		}
	}
	return count
}

func (c *Client) TotalDownloadsByEbook(ebookID uint) int {
	var count int
	for _, purchase := range c.Purchases {
		if purchase.EbookID == ebookID {
			count = +purchase.DownloadsUsed
		}
	}
	return count
}

func (c *Client) GetBirthdateBR() string {
	partsDate := strings.Split(c.Birthdate, "-")
	return fmt.Sprintf("%s/%s/%s", partsDate[2], partsDate[1], partsDate[0])
}
