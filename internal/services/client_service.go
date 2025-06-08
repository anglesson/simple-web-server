package services

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
)

type ClientService struct {
	clientRepository  *repositories.ClientRepository
	creatorRepository *repositories.CreatorRepository
}

func NewClientService(clientRepository *repositories.ClientRepository, creatorRepository *repositories.CreatorRepository) application.ClientServicePort {
	return &ClientService{
		clientRepository:  clientRepository,
		creatorRepository: creatorRepository,
	}
}

func (cs *ClientService) CreateClient(input application.CreateClientInput) (*models.Client, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserEmail(input.EmailCreator)
	if err != nil {
		return nil, err
	}
	client := models.NewClient(input.Name, input.CPF, input.BirthDate, input.Email, input.Phone, creator)

	if err := cs.validateReceita(client); err != nil {
		return nil, err
	}

	err = cs.clientRepository.Save(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cs *ClientService) FindCreatorsClientByID(clientID uint, creatorEmail string) (*models.Client, error) {
	if clientID == 0 || creatorEmail == "" {
		return nil, errors.New("o id do cliente deve ser informado")
	}

	var client models.Client

	err := cs.clientRepository.FindByIDAndCreators(&client, clientID, creatorEmail)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (cs *ClientService) Update(input application.UpdateClientInput) (*models.Client, error) {
	client := &models.Client{}
	if err := cs.validateReceita(client); err != nil {
		return nil, err
	}

	err := cs.clientRepository.Save(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cs *ClientService) CreateBatchClient(clients []*models.Client) error {
	err := cs.clientRepository.InsertBatch(clients)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ClientService) validateReceita(client *models.Client) error {
	rfs := NewReceitaFederalService()
	response := rfs.ConsultaCPF(client.CPF, client.Birthdate)

	if response == nil || !response.Status {
		return errors.New("dados não encontrados na receita federal")
	}

	log.Printf("Erro validateReceita: %v", response)
	client.Name = response.Result.NomeDaPF
	client.Validated = true
	return nil
}
