package client

import (
	"github.com/anglesson/simple-web-server/internal/application/dtos"
	"github.com/anglesson/simple-web-server/internal/models"
)

type ClientServicePort interface {
	CreateClient(input dtos.CreateClientInput) (*models.Client, error)
	FindCreatorsClientByID(clientID uint, creatorEmail string) (*models.Client, error)
	Update(input dtos.UpdateClientInput) (*models.Client, error)
	CreateBatchClient(clients []*models.Client) error
}

type ClientRepositoryPort interface {
	Save(client *models.Client) error
	FindClientsByCreator(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error)
	FindByIDAndCreators(client *models.Client, clientID uint, creator string) error
	FindByClientsWhereEbookNotSend(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error)
	FindByClientsWhereEbookWasSend(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error)
	InsertBatch(clients []*models.Client) error
}

type CreatorRepositoryPort interface {
	FindCreatorByUserID(userID uint) (*models.Creator, error)
	FindCreatorByUserEmail(email string) (*models.Creator, error)
	Create(creator *models.Creator) error
}
