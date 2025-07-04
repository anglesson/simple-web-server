package models

import (
	"github.com/anglesson/simple-web-server/pkg/utils"
	"gorm.io/gorm"
)

type Ebook struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	Status      bool    `json:"status"`
	Image       string  `json:"image"`
	File        string  `json:"file"`
	FileURL     string  `json:"file_url"`
	CreatorID   uint    `json:"creator_id"`
	Creator     Creator `gorm:"foreignKey:CreatorID"`
}

func NewEbook(title, description, file string, value float64, creator Creator) *Ebook {
	return &Ebook{
		Title:       title,
		Description: description,
		Value:       value,
		Status:      true,
		File:        file,
		CreatorID:   creator.ID,
	}
}

func (e *Ebook) GetValue() string {
	return utils.FloatToBRL(e.Value)
}

func (e *Ebook) GetLastUpdate() string {
	return e.UpdatedAt.Format("02-01-2006 15:04")
}
