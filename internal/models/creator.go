package models

import (
	"gorm.io/gorm"
)

type Creator struct {
	*gorm.Model
	Name      string `json:"name"`
	ContactID uint   `json:"contact_id"` // Foreign key
	Contact   Contact
	UserID    uint `json:"user_id"`
	User      User `gorm:"foreignKey:UserID"`
	Ebooks    []Ebook
	Clients   []*Client `gorm:"many2many:client_creator"`
}

func NewCreator(name, email, phone string, user_id uint) *Creator {
	return &Creator{
		Name: name,
		Contact: Contact{
			Email: email,
			Phone: phone,
		},
		UserID: user_id,
	}
}

func (e *Creator) GetEbooks() []Ebook {
	return e.Ebooks
}

func (e *Creator) AddClient(client *Client) {
	e.Clients = append(e.Clients, client)
}
