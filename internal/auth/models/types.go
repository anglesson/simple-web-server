package models

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FormErrors map[string]string

type RegisterForm struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}
