package models

import "gorm.io/gorm"

type Client struct {
	*gorm.Model
	Name      string  `json:"name"`
	CPF       string  `json:"cpf"`
	ContactID uint    `json:"contact_id"`
	Contact   Contact `gorm:"foreignKey:ContactID"`
}
