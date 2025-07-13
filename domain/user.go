package domain

import "errors"

type User struct {
	Entity
	Email    *Email
	Password *Password
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

	passwordVO, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    emailVO,
		Password: passwordVO,
		Username: username,
	}, nil
}

func (u *User) SetPassword(newPassword string) error {
	passwordVO, err := NewPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = passwordVO
	return nil
}
