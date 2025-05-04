package models

import "gorm.io/gorm"

type Ebook struct {
	*gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Value       string  `json:"value"`
	Status      bool    `json:"status"`
	Image       string  `json:"image"`
	File        string  `json:"file"`
	CreatorID   uint    `json:"creator_id"`
	Creator     Creator `gorm:"foreignKey:CreatorID"`
}

func NewEbook(title, description, value, file string, creator Creator) *Ebook {
	return &Ebook{
		Title:       title,
		Description: description,
		Value:       value,
		Status:      true,
		File:        file,
		Creator:     creator,
	}
}
