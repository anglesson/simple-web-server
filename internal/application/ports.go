package application

import "github.com/anglesson/simple-web-server/internal/models"

type ClientServicePort interface {
	CreateClient(input CreateClientInput) (*models.Client, error)
	FindCreatorsClientByID(clientID uint, creatorID uint) (*models.Client, error)
	Update(client *models.Client, input models.ClientRequest) error
	CreateBatchClient(clients []*models.Client) error
}
