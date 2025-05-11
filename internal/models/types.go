package models

type EbookRequest struct {
	Title       string  `validate:"required,min=5,max=120" json:"title"`
	Description string  `validate:"required,max=120" json:"description"`
	Value       float64 `validate:"required,gt=0" json:"value"`
	Status      bool    `validate:"required" json:"status"`
}

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

type ClientRequest struct {
	Name  string `validate:"required,min=5,max=120" json:"name"`
	CPF   string `validate:"required,max=120" json:"cpf"`
	Email string `validate:"required,email" json:"email"`
	Phone string `validate:"max=14" json:"phone"`
}
