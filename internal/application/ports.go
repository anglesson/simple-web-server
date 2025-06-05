package application

import (
	"github.com/anglesson/simple-web-server/internal/domain"
)

type ClientRepositoryInterface interface {
	Create(*domain.Client)
	FindByCPF(cpf string) *domain.Client
}

type ClientUseCasePort interface {
	CreateClient(input CreateClientInput) (*CreateClientOutput, error)
}
