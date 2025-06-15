package client_application

import (
	"errors"

	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	"github.com/anglesson/simple-web-server/internal/models"
)

type ClientService struct {
	clientRepository      ClientRepositoryPort
	creatorRepository     CreatorRepositoryPort
	receitaFederalService common_application.ReceitaFederalServicePort
}

func NewClientService(
	clientRepository ClientRepositoryPort,
	creatorRepository CreatorRepositoryPort,
	receitaFederalService common_application.ReceitaFederalServicePort,
) *ClientService {
	return &ClientService{
		clientRepository:      clientRepository,
		creatorRepository:     creatorRepository,
		receitaFederalService: receitaFederalService,
	}
}

func (cs *ClientService) CreateClient(input CreateClientInput) (*models.Client, error) {
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

func (cs *ClientService) Update(input UpdateClientInput) (*models.Client, error) {
	if input.ID == 0 || input.EmailCreator == "" {
		return nil, errors.New("id do cliente e email do criador são obrigatórios")
	}

	// Find the existing client
	client := &models.Client{}
	err := cs.clientRepository.FindByIDAndCreators(client, input.ID, input.EmailCreator)
	if err != nil {
		return nil, err
	}

	// Update only email and phone
	client.Contact.Email = input.Email
	client.Contact.Phone = input.Phone

	err = cs.clientRepository.Save(client)
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
	if cs.receitaFederalService == nil {
		return errors.New("serviço da receita federal não está disponível")
	}

	response, err := cs.receitaFederalService.ConsultaCPF(client.CPF, client.Birthdate)
	if err != nil {
		return err
	}

	if !response.Status {
		return errors.New("CPF inválido ou não encontrado na receita federal")
	}

	if response.Result.NomeDaPF == "" {
		return errors.New("nome não encontrado na receita federal")
	}

	client.Name = response.Result.NomeDaPF
	client.Validated = true
	return nil
}
