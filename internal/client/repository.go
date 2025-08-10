package client

import (
	"github.com/anglesson/simple-web-server/internal/models"
)

type ClientRepository interface {
	Save(client *models.Client) error
	FindClientsByCreator(creator *models.Creator, query ClientFilter) (*[]models.Client, error)
	FindByIDAndCreators(client *models.Client, clientID uint, creator string) error
	FindByClientsWhereEbookNotSend(creator *models.Creator, query ClientFilter) (*[]models.Client, error)
	FindByClientsWhereEbookWasSend(creator *models.Creator, query ClientFilter) (*[]models.Client, error)
	InsertBatch(clients []*models.Client) error
	FindByEmail(email string) (*models.Client, error)
}
