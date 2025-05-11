package services

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
)

type ClientService struct {
	clientRepository *repositories.ClientRepository
}

func NewClientService() *ClientService {
	r := repositories.NewClientRepository()
	return &ClientService{
		clientRepository: r,
	}
}

func (cs *ClientService) CreateClient(name, cpf, email, phone string, creator *models.Creator) (*models.Client, error) {
	client := models.NewClient(name, cpf, email, phone, creator)
	err := cs.clientRepository.Save(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cs *ClientService) FindCreatorsClientByID(clientID uint, creatorID uint) (*models.Client, error) {
	if clientID == 0 || creatorID == 0 {
		return nil, errors.New("o id do cliente deve ser informado")
	}

	var client models.Client

	err := cs.clientRepository.FindByIDAndCreators(&client, clientID, creatorID)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (cs *ClientService) Update(client *models.Client, input models.ClientRequest) error {
	client.Update(input.Name, input.CPF, input.Email, input.Phone)
	err := cs.clientRepository.Save(client)
	if err != nil {
		return err
	}
	return nil
}
