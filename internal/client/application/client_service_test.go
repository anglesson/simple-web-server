package client_application_test

import (
	"testing"

	client_application "github.com/anglesson/simple-web-server/internal/client/application"
	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var _ client_application.ClientRepositoryPort = (*MockClientRepository)(nil)
var _ client_application.CreatorRepositoryPort = (*MockCreatorRepository)(nil)
var _ common_application.ReceitaFederalServicePort = (*MockRFService)(nil)

type MockCreatorRepository struct {
	mock.Mock
}

func (m *MockCreatorRepository) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) FindCreatorByUserEmail(email string) (*models.Creator, error) {
	args := m.Called(email) // Fixed argument passing
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) Create(creator *models.Creator) error {
	args := m.Called(creator) // Fixed argument passing
	return args.Error(0)      // Fixed error index
}

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Save(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) FindClientsByCreator(creator *models.Creator, query client_application.ClientQuery) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByIDAndCreators(client *models.Client, clientID uint, creator string) error {
	args := m.Called(client, clientID, creator)
	return args.Error(0)
}

func (m *MockClientRepository) FindByClientsWhereEbookNotSend(creator *models.Creator, query client_application.ClientQuery) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByClientsWhereEbookWasSend(creator *models.Creator, query client_application.ClientQuery) (*[]models.Client, error) {
	args := m.Called(creator, query) // Fixed argument passing
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) InsertBatch(clients []*models.Client) error {
	args := m.Called(clients) // Fixed argument passing
	return args.Error(0)
}

type MockRFService struct {
	mock.Mock
}

func (m *MockRFService) ConsultaCPF(cpf, dataNascimento string) (*common_application.ReceitaFederalResponse, error) {
	args := m.Called(cpf, dataNascimento)
	return args.Get(0).(*common_application.ReceitaFederalResponse), args.Error(1)
}

type ClientServiceTestSuite struct {
	suite.Suite
	sut                   *client_application.ClientService
	mockClientRepository  client_application.ClientRepositoryPort
	mockCreatorRepository client_application.CreatorRepositoryPort
	mockRFService         common_application.ReceitaFederalServicePort
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockClientRepository = new(MockClientRepository)
	suite.mockCreatorRepository = new(MockCreatorRepository)
	suite.mockRFService = new(MockRFService)
	suite.sut = client_application.NewClientService(suite.mockClientRepository, suite.mockCreatorRepository, suite.mockRFService)
}

func (suite *ClientServiceTestSuite) TestCreateClient() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}
	expectedClient := &models.Client{
		Name:      "John Doe",
		CPF:       "52998224725",
		Birthdate: "1990-01-01",
		Contact: models.Contact{
			Email: "client@mail.com",
			Phone: "11987654321",
		},
		Creators: []*models.Creator{creator},
	}

	suite.mockCreatorRepository.(*MockCreatorRepository).
		On("FindCreatorByUserEmail", creator.Contact.Email).
		Return(creator, nil)

	suite.mockRFService.(*MockRFService).
		On("ConsultaCPF", expectedClient.CPF, expectedClient.Birthdate).
		Return(&common_application.ReceitaFederalResponse{
			Status: true,
			Result: common_application.ConsultaData{
				NomeDaPF: "John Doe",
			},
		}, nil)

	suite.mockClientRepository.(*MockClientRepository).
		On("Save", mock.MatchedBy(func(client *models.Client) bool {
			return client.Name == expectedClient.Name &&
				client.CPF == expectedClient.CPF &&
				client.Birthdate == expectedClient.Birthdate &&
				client.Contact.Email == expectedClient.Contact.Email &&
				client.Contact.Phone == expectedClient.Contact.Phone &&
				len(client.Creators) == 1 &&
				client.Creators[0] == creator &&
				client.Validated
		})).
		Return(nil)

	input := client_application.CreateClientInput{
		Name:         expectedClient.Name,
		CPF:          expectedClient.CPF,
		BirthDate:    expectedClient.Birthdate,
		Email:        expectedClient.Contact.Email,
		Phone:        expectedClient.Contact.Phone,
		EmailCreator: creator.Contact.Email,
	}

	client, err := suite.sut.CreateClient(input)

	suite.NoError(err)
	suite.NotNil(client)
	suite.Equal(expectedClient.Name, client.Name)
	suite.Equal(expectedClient.CPF, client.CPF)
	suite.Equal(expectedClient.Birthdate, client.Birthdate)
	suite.Equal(expectedClient.Contact.Email, client.Contact.Email)
	suite.Equal(expectedClient.Contact.Phone, client.Contact.Phone)
	suite.Equal(creator, client.Creators[0])
	suite.True(client.Validated)

	suite.mockCreatorRepository.(*MockCreatorRepository).AssertCalled(suite.T(), "FindCreatorByUserEmail", creator.Contact.Email)
	suite.mockRFService.(*MockRFService).AssertCalled(suite.T(), "ConsultaCPF", expectedClient.CPF, expectedClient.Birthdate)
	suite.mockClientRepository.(*MockClientRepository).AssertCalled(suite.T(), "Save", mock.Anything)
}

func (suite *ClientServiceTestSuite) TestUpdateClient() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}
	existingClient := &models.Client{
		Model:     gorm.Model{ID: 1},
		Name:      "John Doe",
		CPF:       "52998224725",
		Birthdate: "1990-01-01",
		Contact: models.Contact{
			Email: "old.email@mail.com",
			Phone: "11987654321",
		},
		Creators:  []*models.Creator{creator},
		Validated: true,
	}

	// Mock finding the existing client
	suite.mockClientRepository.(*MockClientRepository).
		On("FindByIDAndCreators", mock.Anything, existingClient.ID, creator.Contact.Email).
		Return(nil).
		Run(func(args mock.Arguments) {
			client := args.Get(0).(*models.Client)
			*client = *existingClient
		})

	// Mock saving the updated client
	suite.mockClientRepository.(*MockClientRepository).
		On("Save", mock.MatchedBy(func(client *models.Client) bool {
			return client.ID == existingClient.ID &&
				client.Name == existingClient.Name && // Name should remain unchanged
				client.CPF == existingClient.CPF && // CPF should remain unchanged
				client.Birthdate == existingClient.Birthdate && // Birthdate should remain unchanged
				client.Contact.Email == "new.email@mail.com" && // Email should be updated
				client.Contact.Phone == "11999999999" && // Phone should be updated
				len(client.Creators) == 1 &&
				client.Creators[0] == creator &&
				client.Validated
		})).
		Return(nil)

	input := client_application.UpdateClientInput{
		ID:           existingClient.ID,
		Email:        "new.email@mail.com",
		Phone:        "11999999999",
		EmailCreator: creator.Contact.Email,
	}

	client, err := suite.sut.Update(input)

	suite.NoError(err)
	suite.NotNil(client)
	suite.Equal(existingClient.ID, client.ID)
	suite.Equal(existingClient.Name, client.Name)           // Name should remain unchanged
	suite.Equal(existingClient.CPF, client.CPF)             // CPF should remain unchanged
	suite.Equal(existingClient.Birthdate, client.Birthdate) // Birthdate should remain unchanged
	suite.Equal("new.email@mail.com", client.Contact.Email) // Email should be updated
	suite.Equal("11999999999", client.Contact.Phone)        // Phone should be updated
	suite.Equal(creator, client.Creators[0])
	suite.True(client.Validated)

	suite.mockClientRepository.(*MockClientRepository).AssertCalled(suite.T(), "FindByIDAndCreators", mock.Anything, existingClient.ID, creator.Contact.Email)
	suite.mockClientRepository.(*MockClientRepository).AssertCalled(suite.T(), "Save", mock.Anything)
}

func TestClientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
