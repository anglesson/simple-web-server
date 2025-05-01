package models

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}
