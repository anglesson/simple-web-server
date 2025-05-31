package application

type CreateClientInput struct {
	Name     string
	CPF      string
	BirthDay string
	Email    string
	Phone    string
}

type CreateClientOutput struct {
	Name     string
	CPF      string
	BirthDay string
	Email    string
	Phone    string
}
