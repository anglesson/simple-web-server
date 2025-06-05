package application

type CreateClientInput struct {
	Name             string
	CPF              string
	BirthDay         string
	Email            string
	Phone            string
	CreatorUserEmail string
}

type CreateClientOutput struct {
	Name     string
	CPF      string
	BirthDay string
	Email    string
	Phone    string
}
