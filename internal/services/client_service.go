package services

import (
	"errors"
	"log"

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

func (cs *ClientService) CreateClient(name, cpf, dataNascimento, email, phone string, creator *models.Creator) (*models.Client, error) {
	client := models.NewClient(name, cpf, dataNascimento, email, phone, creator)

	if err := cs.validateReceita(client); err != nil {
		return nil, err
	}

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
	if err := cs.validateReceita(client); err != nil {
		return err
	}

	err := cs.clientRepository.Save(client)
	if err != nil {
		return err
	}
	return nil
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
	response := rfs.ConsultaCPF(client.CPF, client.DataNascimento)

	if response == nil || !response.Status {
		return errors.New("dados não encontrados na receita federal")
	}

	log.Printf("Erro validateReceita: %v", response)
	client.Name = response.Result.NomeDaPF
	client.Validated = true
	return nil
}
