package application

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/client/domain"
	"github.com/stretchr/testify/mock"
)

// MockClientRepository is a mock implementation of ClientRepositoryInterface
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(client *domain.Client) {
	m.Called(client)
}

func (m *MockClientRepository) FindByCPF(cpf string) *domain.Client {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.Client)
}

// MockReceitaFederalService is a mock implementation of ReceitaFederalServiceInterface
type MockReceitaFederalService struct {
	mock.Mock
}

func (m *MockReceitaFederalService) Search(cpf, birthDay string) (ReceitaFederalData, error) {
	args := m.Called(cpf, birthDay)
	data := args.Get(0)
	err := args.Get(1)

	if err != nil {
		return data.(ReceitaFederalData), err.(error)
	}
	return data.(ReceitaFederalData), nil
}

func TestCreateClientUseCase_ShoulCallRepositoryWithCorretParam(t *testing.T) {
	input := CreateClientInput{
		Name:     "any_name",
		CPF:      "any_cpf",
		BirthDay: "any_birthday",
		Email:    "any_email",
		Phone:    "any_phone",
	}
	mockRepo := new(MockClientRepository)

	mockRepo.On("FindByCPF", "any_cpf").Return(nil)

	mockRepo.On("Create", &domain.Client{
		Name:     "name_rf",
		CPF:      input.CPF,
		BirthDay: input.BirthDay,
		Email:    input.Email,
		Phone:    input.Phone,
	}).Return(nil)

	mockReceitaService := new(MockReceitaFederalService)
	mockReceitaService.On("Search", "any_cpf", "any_birthday").Return(ReceitaFederalData{NomeDaPF: "name_rf"}, nil)
	createClientUseCase := NewCreateClientUseCase(mockRepo, mockReceitaService)

	_, err := createClientUseCase.Execute(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateClientUseCase_ShouldReturnErrorIfClientAlready(t *testing.T) {
	input := CreateClientInput{
		Name:     "any_name",
		CPF:      "any_cpf",
		BirthDay: "any_birthday",
		Email:    "any_email",
		Phone:    "any_phone",
	}
	mockRepo := new(MockClientRepository)

	mockRepo.On("FindByCPF", "any_cpf").Return(&domain.Client{})

	mockReceitaService := new(MockReceitaFederalService)
	mockReceitaService.On("Search", "any_cpf", "any_birthday")
	clientUseCase := NewCreateClientUseCase(mockRepo, mockReceitaService)

	_, err := clientUseCase.Execute(input)

	if err == nil {
		t.Error("Expected error when client already exists, got nil")
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateClientUseCase_ShouldCallReceitaFederalService(t *testing.T) {
	input := CreateClientInput{
		Name:     "any_name",
		CPF:      "any_cpf",
		BirthDay: "any_birthday",
		Email:    "any_email",
		Phone:    "any_phone",
	}
	mockRepo := new(MockClientRepository)
	mockRepo.On("FindByCPF", "any_cpf").Return(nil)
	mockRepo.On("Create", &domain.Client{
		Name:     "name_pf",
		CPF:      input.CPF,
		BirthDay: input.BirthDay,
		Email:    input.Email,
		Phone:    input.Phone,
	}).Return(nil)

	mockRF := new(MockReceitaFederalService)
	mockRF.On("Search", "any_cpf", "any_birthday").Return(ReceitaFederalData{
		NomeDaPF: "name_pf",
	}, nil)

	createClientUseCase := NewCreateClientUseCase(mockRepo, mockRF)

	_, err := createClientUseCase.Execute(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	mockRF.AssertExpectations(t)
}
