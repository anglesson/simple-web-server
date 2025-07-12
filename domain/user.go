package domain

import "errors"

type User struct {
	Email    *Email
	Password string
	Username string
}

func NewUser(username, email, password string) (*User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email and password are required")
	}

	if len(username) > 50 {
		username = username[0:50]
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    emailVO,
		Password: password,
		Username: username,
	}, nil
}
