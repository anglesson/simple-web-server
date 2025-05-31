package application

import "github.com/anglesson/simple-web-server/internal/client/domain"

type ClientRepositoryInterface interface {
	Create(*domain.Client)
	FindByCPF(cpf string) *domain.Client
}
