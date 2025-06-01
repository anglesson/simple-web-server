package client_test

import (
	"testing"

	application "github.com/anglesson/simple-web-server/internal/application/client"
	common "github.com/anglesson/simple-web-server/internal/application/common"
	domain "github.com/anglesson/simple-web-server/internal/domain/client"
	"github.com/stretchr/testify/assert"
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

func (m *MockReceitaFederalService) Search(cpf, birthDay string) (common.ReceitaFederalData, error) {
	args := m.Called(cpf, birthDay)
	data := args.Get(0)
	err := args.Get(1)

	if err != nil {
		return data.(common.ReceitaFederalData), err.(error)
	}
	return data.(common.ReceitaFederalData), nil
}

type testSetup struct {
	mockRepo       *MockClientRepository
	mockRFService  *MockReceitaFederalService
	createClientUC *application.CreateClientUseCase
	defaultInput   application.CreateClientInput
}

func setupTest(t *testing.T) *testSetup {
	mockRepo := new(MockClientRepository)
	mockRFService := new(MockReceitaFederalService)
	createClientUC := application.NewCreateClientUseCase(mockRepo, mockRFService)

	defaultInput := application.CreateClientInput{
		Name:     "any_name",
		CPF:      "any_cpf",
		BirthDay: "any_birthday",
		Email:    "any_email",
		Phone:    "any_phone",
	}

	return &testSetup{
		mockRepo:       mockRepo,
		mockRFService:  mockRFService,
		createClientUC: createClientUC,
		defaultInput:   defaultInput,
	}
}

func TestCreateClientUseCase(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*testSetup)
		expectedError bool
		errorMessage  string
	}{
		{
			name: "should create client successfully",
			setupMocks: func(ts *testSetup) {
				ts.mockRepo.On("FindByCPF", "any_cpf").Return(nil)
				ts.mockRFService.On("Search", "any_cpf", "any_birthday").Return(common.ReceitaFederalData{NomeDaPF: "name_rf"}, nil)
				ts.mockRepo.On("Create", &domain.Client{
					Name:     "name_rf",
					CPF:      ts.defaultInput.CPF,
					BirthDay: ts.defaultInput.BirthDay,
					Email:    ts.defaultInput.Email,
					Phone:    ts.defaultInput.Phone,
				}).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should return error if client already exists",
			setupMocks: func(ts *testSetup) {
				ts.mockRepo.On("FindByCPF", "any_cpf").Return(&domain.Client{})
			},
			expectedError: true,
			errorMessage:  "client already exists",
		},
		{
			name: "should return error if Receita Federal service fails",
			setupMocks: func(ts *testSetup) {
				ts.mockRepo.On("FindByCPF", "any_cpf").Return(nil)
				ts.mockRFService.On("Search", "any_cpf", "any_birthday").Return(common.ReceitaFederalData{}, assert.AnError)
			},
			expectedError: true,
			errorMessage:  "failed to validate CPF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)
			tt.setupMocks(ts)

			_, err := ts.createClientUC.Execute(ts.defaultInput)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}

			ts.mockRepo.AssertExpectations(t)
			ts.mockRFService.AssertExpectations(t)
		})
	}
}
