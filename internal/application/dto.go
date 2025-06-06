package application

import domain "github.com/anglesson/simple-web-server/internal/domain"

type CreateClientInput struct {
	Name             string
	CPF              string
	BirthDay         string
	Email            string
	Phone            string
	CreatorUserEmail string
}

type CreateClientOutput struct {
	ID uint
}

type UpdateClientInput struct {
	ID               uint
	Name             string
	CPF              string
	BirthDay         string
	Email            string
	Phone            string
	CreatorUserEmail string
}

type UpdateClientOutput struct {
	ID uint
}

type ImportClientsInput struct {
	File             []byte
	FileName         string
	CreatorUserEmail string
}

type ImportClientsOutput struct {
	ImportedCount int
}

type ListClientsInput struct {
	Term             string
	Page             int
	PageSize         int
	CreatorUserEmail string
}

type ListClientsOutput struct {
	Clients    []*domain.Client
	TotalCount int
	Page       int
	PageSize   int
}

type Pagination struct {
	Page     int
	PageSize int
	Total    int
}
