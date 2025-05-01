package models

import "gorm.io/gorm"

type Contact struct {
	*gorm.Model
	Email string `json:"email"`
	Phone string `json:"phone"`
}
