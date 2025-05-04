package models

import "gorm.io/gorm"

type Purchase struct {
	*gorm.Model
	EbookID  uint   `json:"ebook_id"`
	Ebook    Ebook  `gorm:"foreignKey:EbookID"`
	ClientID uint   `json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID"`
}
