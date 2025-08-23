package data

import (
	"gorm.io/gorm"
)

type AccountModel struct {
	gorm.Model
	Origin    string
	AccountID string
	SellerID  uint
	IsPending bool
}
