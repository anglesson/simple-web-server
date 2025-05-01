package models

import "gorm.io/gorm"

type Purchase struct {
	*gorm.Model
	Ebook  Ebook
	Client Client
}
