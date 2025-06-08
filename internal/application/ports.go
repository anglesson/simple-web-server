package application

import "github.com/anglesson/simple-web-server/internal/models"

type ClientServicePort interface {
	CreateClient(input CreateClientInput) (*models.Client, error)
	FindCreatorsClientByID(clientID uint, creatorID uint) (*models.Client, error)
	Update(input UpdateClientInput) (*models.Client, error)
	CreateBatchClient(clients []*models.Client) error
}
