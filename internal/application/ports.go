package application

import (
	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/anglesson/simple-web-server/internal/domain"
)

type ClientRepositoryInterface interface {
	Create(client *domain.Client) error
	Update(client *domain.Client) error
	FindByCPF(cpf common_domain.CPF) *domain.Client
	FindByID(id uint) (*domain.Client, error)
	FindByCreatorID(creatorID uint, query ClientQuery) ([]*domain.Client, error)
	CreateBatch(clients []*domain.Client) error
	List(query ClientQuery) ([]*domain.Client, int, error)
}

type ClientUseCasePort interface {
	CreateClient(input CreateClientInput) (*CreateClientOutput, error)
	UpdateClient(input UpdateClientInput) (*UpdateClientOutput, error)
	ImportClients(input ImportClientsInput) (*ImportClientsOutput, error)
	ListClients(input ListClientsInput) (*ListClientsOutput, error)
}

type ClientQuery struct {
	Term       string
	Pagination *Pagination
	CreatorID  uint
}
