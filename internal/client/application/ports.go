package client_application

import (
	client_domain "github.com/anglesson/simple-web-server/internal/client/domain"
)

type ClientRepositoryInterface interface {
	Create(*client_domain.Client)
	FindByCPF(cpf string) *client_domain.Client
}
