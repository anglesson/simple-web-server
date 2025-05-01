package models

type Client struct {
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	Contact Contact
}
