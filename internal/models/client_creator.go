package models

import "gorm.io/gorm"

type ClientCreator struct {
	gorm.Model
	ClientID  uint    `json:"client_id"`
	Client    Client  `gorm:"foreignKey:ClientID"`
	CreatorID uint    `json:"creator_id"`
	Creator   Creator `gorm:"foreignKey:CreatorID"`
}
