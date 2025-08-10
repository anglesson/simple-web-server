package client

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
)

// Module representa um módulo completo com todas as suas dependências
type Module struct {
	Handler    *ClientHandler
	Service    service.ClientService
	Repository repository.ClientRepository
}

// ClientFilter é usado para filtrar clientes
type ClientFilter struct {
	Term       string
	EbookID    uint
	Pagination *models.Pagination
}
