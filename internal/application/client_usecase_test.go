package application_test

import (
	"fmt"
	"testing"

	application "github.com/anglesson/simple-web-server/internal/application"
	domain "github.com/anglesson/simple-web-server/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClientRepository is a mock implementation of ClientRepositoryInterface
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(client *domain.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) Update(client *domain.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) FindByCPF(cpf domain.CPF) *domain.Client {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.Client)
}

func (m *MockClientRepository) FindByID(id uint) (*domain.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockClientRepository) FindByCreatorID(creatorID uint, query application.ClientQuery) ([]*domain.Client, error) {
	args := m.Called(creatorID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Client), args.Error(1)
}

func (m *MockClientRepository) CreateBatch(clients []*domain.Client) error {
	args := m.Called(clients)
	return args.Error(0)
}

func (m *MockClientRepository) List(query application.ClientQuery) ([]*domain.Client, int, error) {
	args := m.Called(query)
	return args.Get(0).([]*domain.Client), args.Get(1).(int), args.Error(2)
}

// CPFService is a mock implementation of CPFServicePort
type CPFService struct {
	mock.Mock
}

func (m *CPFService) ConsultCPF(cpf domain.CPF, birthDay domain.BirthDate) (application.CPFOutput, error) {
	args := m.Called(cpf, birthDay)
	return args.Get(0).(application.CPFOutput), args.Error(1)
}

type testSetup struct {
	mockRepo       *MockClientRepository
	mockRFService  *CPFService
	createClientUC *application.ClientUseCase
	defaultInput   application.CreateClientInput
}

func setupTest(t *testing.T) *testSetup {
	mockRepo := new(MockClientRepository)
	mockRFService := new(CPFService)
	createClientUC := application.NewClientUseCase(mockRepo, mockRFService)

	defaultInput := application.CreateClientInput{
		Name:     "any_name",
		CPF:      "461.371.640-39",
		BirthDay: "1990-01-01",
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
				cpf, _ := domain.NewCPF(ts.defaultInput.CPF)
				ts.mockRepo.On("FindByCPF", cpf).Return(nil)
				birthDay, _ := domain.NewBirthDate(ts.defaultInput.BirthDay)
				ts.mockRFService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{NomeDaPF: "name_rf"}, nil)
				ts.mockRepo.On("Create", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should return error if client already exists",
			setupMocks: func(ts *testSetup) {
				cpf, _ := domain.NewCPF(ts.defaultInput.CPF)
				ts.mockRepo.On("FindByCPF", cpf).Return(&domain.Client{})
			},
			expectedError: true,
			errorMessage:  "client already exists",
		},
		{
			name: "should return error if Receita Federal service fails",
			setupMocks: func(ts *testSetup) {
				cpf, _ := domain.NewCPF(ts.defaultInput.CPF)
				ts.mockRepo.On("FindByCPF", cpf).Return(nil)
				birthDay, _ := domain.NewBirthDate(ts.defaultInput.BirthDay)
				ts.mockRFService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{}, assert.AnError)
			},
			expectedError: true,
			errorMessage:  "failed to validate CPF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setupTest(t)
			tt.setupMocks(ts)

			_, err := ts.createClientUC.CreateClient(ts.defaultInput)

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

func TestClientUseCase_CreateClient(t *testing.T) {
	tests := []struct {
		name          string
		input         application.CreateClientInput
		mockSetup     func(*MockClientRepository, *CPFService)
		expectedError bool
	}{
		{
			name: "successful creation",
			input: application.CreateClientInput{
				Name:             "John Doe",
				CPF:              "461.371.640-39",
				BirthDay:         "1990-01-01",
				Email:            "john@example.com",
				Phone:            "1234567890",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				cpf, _ := domain.NewCPF("461.371.640-39")
				repo.On("FindByCPF", cpf).Return(nil)
				birthDay, _ := domain.NewBirthDate("1990-01-01")
				cpfService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{
					NomeDaPF: "John Doe",
				}, nil)
				repo.On("Create", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "client already exists",
			input: application.CreateClientInput{
				Name:             "John Doe",
				CPF:              "461.371.640-39",
				BirthDay:         "1990-01-01",
				Email:            "john@example.com",
				Phone:            "1234567890",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				cpf, _ := domain.NewCPF("461.371.640-39")
				repo.On("FindByCPF", cpf).Return(&domain.Client{})
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockClientRepository)
			mockCPFService := new(CPFService)
			tt.mockSetup(mockRepo, mockCPFService)

			useCase := application.NewClientUseCase(mockRepo, mockCPFService)
			_, err := useCase.CreateClient(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("CreateClient() error = %v, wantErr %v", err, tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
			mockCPFService.AssertExpectations(t)
		})
	}
}

func TestClientUseCase_UpdateClient(t *testing.T) {
	tests := []struct {
		name          string
		input         application.UpdateClientInput
		mockSetup     func(*MockClientRepository, *CPFService)
		expectedError bool
	}{
		{
			name: "successful update",
			input: application.UpdateClientInput{
				ID:               1,
				Name:             "John Doe",
				CPF:              "461.371.640-39",
				BirthDay:         "1990-01-01",
				Email:            "john@example.com",
				Phone:            "1234567890",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				client, _ := domain.NewClient("Old Name", "461.371.640-39", "1990-01-01", "old@example.com", "1234567890")
				repo.On("FindByID", uint(1)).Return(client, nil)
				cpf, _ := domain.NewCPF("461.371.640-39")
				birthDay, _ := domain.NewBirthDate("1990-01-01")
				cpfService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{
					NomeDaPF: "John Doe",
				}, nil)
				repo.On("Update", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "client not found",
			input: application.UpdateClientInput{
				ID:               1,
				Name:             "John Doe",
				CPF:              "461.371.640-39",
				BirthDay:         "1990-01-01",
				Email:            "john@example.com",
				Phone:            "1234567890",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				repo.On("FindByID", uint(1)).Return(nil, fmt.Errorf("client not found"))
			},
			expectedError: true,
		},
		{
			name: "receita federal validation fails",
			input: application.UpdateClientInput{
				ID:               1,
				Name:             "John Doe",
				CPF:              "461.371.640-39",
				BirthDay:         "1990-01-01",
				Email:            "john@example.com",
				Phone:            "1234567890",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				client, _ := domain.NewClient("Old Name", "461.371.640-39", "1990-01-01", "old@example.com", "1234567890")
				repo.On("FindByID", uint(1)).Return(client, nil)
				cpf, _ := domain.NewCPF("461.371.640-39")
				birthDay, _ := domain.NewBirthDate("1990-01-01")
				cpfService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{}, fmt.Errorf("receita federal validation failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockClientRepository)
			mockCPFService := new(CPFService)
			tt.mockSetup(mockRepo, mockCPFService)

			useCase := application.NewClientUseCase(mockRepo, mockCPFService)
			_, err := useCase.UpdateClient(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("UpdateClient() error = %v, wantErr %v", err, tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
			mockCPFService.AssertExpectations(t)
		})
	}
}

func TestClientUseCase_ImportClients(t *testing.T) {
	tests := []struct {
		name          string
		input         application.ImportClientsInput
		mockSetup     func(*MockClientRepository, *CPFService)
		expectedError bool
	}{
		{
			name: "successful import",
			input: application.ImportClientsInput{
				File:             []byte("name,cpf,birth_day,email,phone\nJohn Doe,461.371.640-39,1990-01-01,john@example.com,1234567890"),
				FileName:         "clients.csv",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				cpf, _ := domain.NewCPF("461.371.640-39")
				birthDay, _ := domain.NewBirthDate("1990-01-01")
				cpfService.On("ConsultCPF", cpf, *birthDay).Return(application.CPFOutput{
					NomeDaPF: "John Doe",
				}, nil)
				repo.On("CreateBatch", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "invalid file format",
			input: application.ImportClientsInput{
				File:             []byte("invalid,csv,format"),
				FileName:         "clients.txt",
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				// No mock setup needed
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockClientRepository)
			mockCPFService := new(CPFService)
			tt.mockSetup(mockRepo, mockCPFService)

			useCase := application.NewClientUseCase(mockRepo, mockCPFService)
			_, err := useCase.ImportClients(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("ImportClients() error = %v, wantErr %v", err, tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
			mockCPFService.AssertExpectations(t)
		})
	}
}

func TestClientUseCase_ListClients(t *testing.T) {
	tests := []struct {
		name          string
		input         application.ListClientsInput
		mockSetup     func(*MockClientRepository, *CPFService)
		expectedError bool
	}{
		{
			name: "successful list",
			input: application.ListClientsInput{
				Term:             "John",
				Page:             1,
				PageSize:         10,
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				client, _ := domain.NewClient("John Doe", "461.371.640-39", "1990-01-01", "john@example.com", "1234567890")
				repo.On("List", mock.Anything).Return([]*domain.Client{client}, 1, nil)
			},
			expectedError: false,
		},
		{
			name: "empty list",
			input: application.ListClientsInput{
				Term:             "",
				Page:             1,
				PageSize:         10,
				CreatorUserEmail: "creator@example.com",
			},
			mockSetup: func(repo *MockClientRepository, cpfService *CPFService) {
				repo.On("List", mock.Anything).Return([]*domain.Client{}, 0, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockClientRepository)
			mockCPFService := new(CPFService)
			tt.mockSetup(mockRepo, mockCPFService)

			useCase := application.NewClientUseCase(mockRepo, mockCPFService)
			_, err := useCase.ListClients(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("ListClients() error = %v, wantErr %v", err, tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
			mockCPFService.AssertExpectations(t)
		})
	}
}
