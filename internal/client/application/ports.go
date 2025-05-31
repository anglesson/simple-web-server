package application

import "github.com/anglesson/simple-web-server/internal/client/domain"

type ClientRepositoryInterface interface {
	Create(*domain.Client)
	FindByCPF(cpf string) *domain.Client
}

type ReceitaFederalServiceInterface interface {
	Search(cpf, birthDay string) (any, error)
}
