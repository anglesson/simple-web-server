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

func TestShoulCallRepositoryWithCorretParam(t *testing.T) {
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
		Name:     input.Name,
		CPF:      input.CPF,
		BirthDay: input.BirthDay,
		Email:    input.Email,
		Phone:    input.Phone,
	}).Return(nil)

	createClientUseCase := NewCreateClientUseCase(mockRepo)

	_, err := createClientUseCase.Execute(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestShouldReturnErrorIfClientAlready(t *testing.T) {
	input := CreateClientInput{
		Name:     "any_name",
		CPF:      "any_cpf",
		BirthDay: "any_birthday",
		Email:    "any_email",
		Phone:    "any_phone",
	}
	mockRepo := new(MockClientRepository)

	mockRepo.On("FindByCPF", "any_cpf").Return(&domain.Client{})

	clientUseCase := NewCreateClientUseCase(mockRepo)

	_, err := clientUseCase.Execute(input)

	if err == nil {
		t.Error("Expected error when client already exists, got nil")
	}

	mockRepo.AssertExpectations(t)
}
