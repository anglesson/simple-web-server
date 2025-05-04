package models

import (
	"github.com/anglesson/simple-web-server/internal/auth/models"
	"gorm.io/gorm"
)

type Creator struct {
	*gorm.Model
	Name      string `json:"name"`
	ContactID uint   `json:"contact_id"` // Foreign key
	Contact   Contact
	UserID    uint        `json:"user_id"`
	User      models.User `gorm:"foreignKey:UserID"`
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
