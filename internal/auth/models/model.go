package models

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}
type User struct {
	ID       int
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
