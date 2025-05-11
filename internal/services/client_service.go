package services

import (
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
