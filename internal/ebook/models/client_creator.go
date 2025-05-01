package models

import "gorm.io/gorm"

type ClientCreator struct {
	*gorm.Model
	Client  Client
	Creator Creator
}
