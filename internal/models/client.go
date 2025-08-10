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
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Validated bool       `json:"validated"`
	Creators  []*Creator `gorm:"many2many:client_creators"`
	Purchases []*Purchase
}

func NewClient(name, cpf, birthDate, email, phone string, creator *Creator) *Client {
	return &Client{
		Name:      name,
		CPF:       cpf,
		Birthdate: birthDate,
		Email:     email,
		Phone:     phone,
		Validated: false,
		Creators: []*Creator{
			creator,
		},
	}
}

func (c *Client) Update(name, cpf, email, phone string) {
	c.Name = name
	c.CPF = cpf
	c.Email = email
	c.Phone = phone
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
	if c.Birthdate == "" {
		return "N/A"
	}

	partsDate := strings.Split(c.Birthdate, "-")
	if len(partsDate) != 3 {
		return "Data inválida"
	}

	return fmt.Sprintf("%s/%s/%s", partsDate[2], partsDate[1], partsDate[0])
}

func (c *Client) GetInitials() string {
	if c.Name == "" {
		return "?"
	}

	names := strings.Fields(c.Name)
	if len(names) == 0 {
		return "?"
	}

	initials := ""
	if len(names) >= 1 && len(names[0]) > 0 {
		initials += string(names[0][0])
	}
	if len(names) >= 2 && len(names[len(names)-1]) > 0 {
		initials += string(names[len(names)-1][0])
	}

	if initials == "" {
		return "?"
	}

	return strings.ToUpper(initials)
}
