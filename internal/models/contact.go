package models

import "gorm.io/gorm"

type Contact struct {
	*gorm.Model
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func NewContact(email, phone string) Contact {
	return Contact{
		Email: email,
		Phone: phone,
	}
}
