package services

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/models"
)

type ClientService interface {
	CreateClient(input CreateClientInput) (*CreateClientOutput, error)
	FindCreatorsClientByID(clientID uint, creatorEmail string) (*models.Client, error)
	Update(input UpdateClientInput) (*models.Client, error)
	CreateBatchClient(clients []*models.Client) error
}

type CreateClientInput struct {
	Name         string
	CPF          string
	Phone        string
	BirthDate    string
	Email        string
	EmailCreator string
}

type CreateClientOutput struct {
	ID        int
	Name      string
	CPF       string
	Phone     string
	BirthDate string
	Email     string
	CreatedAt string
	UpdatedAt string
}

type UpdateClientInput struct {
	ID           uint
	Email        string
	Phone        string
	EmailCreator string
}

type clientServiceImpl struct {
	clientRepository      repository.ClientRepository
	creatorRepository     repository.CreatorRepository
	receitaFederalService gov.ReceitaFederalService
}

func NewClientService(
	clientRepository repository.ClientRepository,
	creatorRepository repository.CreatorRepository,
	receitaFederalService gov.ReceitaFederalService,
) ClientService {
	return &clientServiceImpl{
		clientRepository:      clientRepository,
		creatorRepository:     creatorRepository,
		receitaFederalService: receitaFederalService,
	}
}

func (cs *clientServiceImpl) CreateClient(input CreateClientInput) (*CreateClientOutput, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserEmail(input.EmailCreator)
	if err != nil {
		return nil, err
	}

	clientExists, err := cs.clientRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if clientExists != nil {
		return nil, errors.New("cliente já existe")
	}

	client := models.NewClient(input.Name, input.CPF, input.BirthDate, input.Email, input.Phone, creator)

	if err := cs.validateReceita(client); err != nil {
		return nil, err
	}

	err = cs.clientRepository.Save(client)
	if err != nil {
		return nil, err
	}
	return &CreateClientOutput{
		ID:        int(client.ID),
		Name:      client.Name,
		CPF:       client.CPF,
		Phone:     client.Contact.Phone,
		Email:     client.Contact.Email,
		BirthDate: client.Birthdate,
		UpdatedAt: client.UpdatedAt.String(),
		CreatedAt: client.CreatedAt.String(),
	}, nil
}

func (cs *clientServiceImpl) FindCreatorsClientByID(clientID uint, creatorEmail string) (*models.Client, error) {
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

func (cs *clientServiceImpl) Update(input UpdateClientInput) (*models.Client, error) {
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

func (cs *clientServiceImpl) CreateBatchClient(clients []*models.Client) error {
	err := cs.clientRepository.InsertBatch(clients)
	if err != nil {
		return err
	}
	return nil
}

func (cs *clientServiceImpl) validateReceita(client *models.Client) error {
	if cs.receitaFederalService == nil {
		return errors.New("serviço da receita federal não está disponível")
	}

	response, err := cs.receitaFederalService.ConsultaCPF(client.CPF, client.GetBirthdateBR())
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
