package dtos

type CreateClientInput struct {
	Name         string
	CPF          string
	Phone        string
	BirthDate    string
	Email        string
	EmailCreator string
}

type UpdateClientInput struct {
	CPF          string
	Phone        string
	Email        string
	EmailCreator string
}

type ClientQuery struct {
	Term       string
	EbookID    uint
	Pagination *Pagination
}
