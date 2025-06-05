package client_application

import (
	"testing"

	client_domain "github.com/anglesson/simple-web-server/internal/client/domain"
	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClientRepository is a mock common_application. of ClientRepositoryInterface
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(client *client_domain.Client) {
	m.Called(client)
}

func (m *MockClientRepository) FindByCPF(cpf string) *client_domain.Client {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*client_domain.Client)
}

// CPFService is a mock implementation of ReceitaFederalServiceInterface
type CPFService struct {
	mock.Mock
}

func (m *CPFService) ConsultCPF(cpf common_domain.CPF, birthDay common_domain.BirthDate) (common_application.CPFOutput, error) {
	args := m.Called(cpf, birthDay)
	return args.Get(0).(common_application.CPFOutput), args.Error(1)
}

type testSetup struct {
	mockRepo       *MockClientRepository
	mockRFService  *CPFService
	createClientUC *ClientUseCase
	defaultInput   CreateClientInput
}

func setupTest(t *testing.T) *testSetup {
	mockRepo := new(MockClientRepository)
	mockRFService := new(CPFService)
	createClientUC := NewClientUseCase(mockRepo, mockRFService)

	defaultInput := CreateClientInput{
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
				ts.mockRepo.On("FindByCPF", "461.371.640-39").Return(nil)
				cpf, _ := common_domain.NewCPF(ts.defaultInput.CPF)
				birthDay, _ := common_domain.NewBirthDate(ts.defaultInput.BirthDay)
				ts.mockRFService.On("ConsultCPF", cpf, *birthDay).Return(common_application.CPFOutput{NomeDaPF: "name_rf"}, nil)
				ts.mockRepo.On("Create", &client_domain.Client{
					Name:     "name_rf",
					CPF:      cpf,
					BirthDay: birthDay,
					Email:    ts.defaultInput.Email,
					Phone:    ts.defaultInput.Phone,
				}).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should return error if client already exists",
			setupMocks: func(ts *testSetup) {
				ts.mockRepo.On("FindByCPF", "461.371.640-39").Return(&client_domain.Client{})
			},
			expectedError: true,
			errorMessage:  "client already exists",
		},
		{
			name: "should return error if Receita Federal service fails",
			setupMocks: func(ts *testSetup) {
				ts.mockRepo.On("FindByCPF", "461.371.640-39").Return(nil)
				cpf, _ := common_domain.NewCPF(ts.defaultInput.CPF)
				birthDay, _ := common_domain.NewBirthDate(ts.defaultInput.BirthDay)
				ts.mockRFService.On("ConsultCPF", cpf, *birthDay).Return(common_application.CPFOutput{}, assert.AnError)
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
