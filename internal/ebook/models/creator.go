package models

import "gorm.io/gorm"

type Creator struct {
	*gorm.Model
	Name    string `json:"name"`
	Contact Contact
}
