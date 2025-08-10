package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/gov"
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
	// Validate input
	if err := validateClientInput(input); err != nil {
		return nil, err
	}

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

	// Clean CPF and phone
	cleanCPF := cleanCPF(input.CPF)
	cleanPhone := cleanPhone(input.Phone)

	// Parse birth date - try DD/MM/YYYY format first (from jmask), then YYYY-MM-DD
	var birthDate time.Time
	birthDate, err = time.Parse("02/01/2006", input.BirthDate)
	if err != nil {
		// If that fails, try YYYY-MM-DD format (from HTML date input)
		birthDate, err = time.Parse("2006-01-02", input.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("formato de data de nascimento inválido: %w", err)
		}
	}

	client := models.NewClient(input.Name, cleanCPF, input.BirthDate, input.Email, cleanPhone, creator)

	// Validate CPF with Receita Federal
	if err := cs.validateReceita(client, birthDate); err != nil {
		return nil, err
	}

	if err := cs.clientRepository.Save(client); err != nil {
		return nil, err
	}

	return &CreateClientOutput{
		ID:        int(client.ID),
		Name:      client.Name,
		CPF:       client.CPF,
		Phone:     client.Phone,
		BirthDate: client.Birthdate,
		Email:     client.Email,
		CreatedAt: client.CreatedAt.Format("02/01/2006 15:04:05"),
		UpdatedAt: client.UpdatedAt.Format("02/01/2006 15:04:05"),
	}, nil
}

func (cs *clientServiceImpl) FindCreatorsClientByID(clientID uint, creatorEmail string) (*models.Client, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserEmail(creatorEmail)
	if err != nil {
		return nil, err
	}

	client := &models.Client{}
	if err := cs.clientRepository.FindByIDAndCreators(client, clientID, creator.Email); err != nil {
		return nil, err
	}

	return client, nil
}

func (cs *clientServiceImpl) Update(input UpdateClientInput) (*models.Client, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserEmail(input.EmailCreator)
	if err != nil {
		return nil, err
	}

	client := &models.Client{}
	if err := cs.clientRepository.FindByIDAndCreators(client, input.ID, creator.Email); err != nil {
		return nil, err
	}

	client.Email = input.Email
	client.Phone = cleanPhone(input.Phone)

	if err := cs.clientRepository.Save(client); err != nil {
		return nil, err
	}

	return client, nil
}

func (cs *clientServiceImpl) CreateBatchClient(clients []*models.Client) error {
	if len(clients) == 0 {
		return errors.New("nenhum cliente para criar")
	}

	return cs.clientRepository.InsertBatch(clients)
}

func (cs *clientServiceImpl) validateReceita(client *models.Client, birthDate time.Time) error {
	// Validate CPF with Receita Federal
	receitaResponse, err := cs.receitaFederalService.ConsultaCPF(client.CPF, birthDate.Format("02/01/2006"))
	if err != nil {
		return fmt.Errorf("erro ao validar CPF na Receita Federal: %w", err)
	}

	if !receitaResponse.Status {
		return fmt.Errorf("CPF inválido na Receita Federal: %s", receitaResponse.Return)
	}

	return nil
}

// Helper functions
func validateClientInput(input CreateClientInput) error {
	if input.Name == "" {
		return errors.New("nome é obrigatório")
	}
	if input.CPF == "" {
		return errors.New("CPF é obrigatório")
	}
	if input.Email == "" {
		return errors.New("email é obrigatório")
	}
	if input.BirthDate == "" {
		return errors.New("data de nascimento é obrigatória")
	}
	if input.EmailCreator == "" {
		return errors.New("email do criador é obrigatório")
	}
	return nil
}

func cleanCPF(cpf string) string {
	// Remove all non-digit characters
	cleaned := ""
	for _, char := range cpf {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}
	return cleaned
}

func cleanPhone(phone string) string {
	// Remove all non-digit characters
	cleaned := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}
	return cleaned
}
