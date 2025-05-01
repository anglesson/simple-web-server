package models

import "gorm.io/gorm"

type Ebook struct {
	*gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Status      string `json:"status"`
	Image       string `json:"image"`
}
