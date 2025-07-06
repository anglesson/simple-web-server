package services_test

import (
	mocks_repo "github.com/anglesson/simple-web-server/internal/repositories/mocks"
	"github.com/anglesson/simple-web-server/pkg/gov/mocks"
	"testing"

	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var _ repositories.ClientRepository = (*MockClientRepository)(nil)
var _ repositories.CreatorRepository = (*mocks_repo.MockCreatorRepository)(nil)
var _ gov.ReceitaFederalService = (*mocks.MockRFService)(nil)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Save(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) FindClientsByCreator(creator *models.Creator, query domain.ClientFilter) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByIDAndCreators(client *models.Client, clientID uint, creator string) error {
	args := m.Called(client, clientID, creator)
	return args.Error(0)
}

func (m *MockClientRepository) FindByClientsWhereEbookNotSend(creator *models.Creator, query domain.ClientFilter) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByClientsWhereEbookWasSend(creator *models.Creator, query domain.ClientFilter) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) InsertBatch(clients []*models.Client) error {
	args := m.Called(clients)
	return args.Error(0)
}

func (m *MockClientRepository) FindByEmail(email string) (*models.Client, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

type ClientServiceTestSuite struct {
	suite.Suite
	sut                   services.ClientService
	mockClientRepository  repositories.ClientRepository
	mockCreatorRepository repositories.CreatorRepository
	mockRFService         gov.ReceitaFederalService
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockClientRepository = new(MockClientRepository)
	suite.mockCreatorRepository = new(mocks_repo.MockCreatorRepository)
	suite.mockRFService = new(mocks.MockRFService)
	suite.sut = services.NewClientService(suite.mockClientRepository, suite.mockCreatorRepository, suite.mockRFService)
}

func (suite *ClientServiceTestSuite) TestCreateClient() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}

	input := services.CreateClientInput{
		Name:         "Name User",
		CPF:          "000.000.000-00",
		BirthDate:    "2012-12-12",
		EmailCreator: creator.Contact.Email,
	}

	expectedName := "Name Receita Federal"
	expectedBirthDay := "12/12/2012"

	client := &models.Client{
		Validated: true,
		Name:      expectedName,
		CPF:       input.CPF,
		Birthdate: input.BirthDate,
		Creators:  []*models.Creator{creator},
	}

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", "000.000.000-00", "12/12/2012").
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedName,
				NumeroDeCPF:    input.CPF,
				DataNascimento: expectedBirthDay,
			},
		}, nil)

	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).
		On("FindCreatorByUserEmail", creator.Contact.Email).
		Return(creator, nil)

	suite.mockClientRepository.(*MockClientRepository).
		On("FindByEmail", input.Email).
		Return(nil, nil)

	suite.mockClientRepository.(*MockClientRepository).
		On("Save", client).
		Return(nil)

	_, err := suite.sut.CreateClient(input)

	suite.NoError(err)
	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindCreatorByUserEmail", creator.Contact.Email)
	suite.mockClientRepository.(*MockClientRepository).AssertCalled(suite.T(), "Save", client)
}

func (suite *ClientServiceTestSuite) TestShouldReturnErrorIfClientExists() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}

	input := services.CreateClientInput{
		Name:         "Name User",
		CPF:          "000.000.000-00",
		BirthDate:    "2012-12-12",
		EmailCreator: creator.Contact.Email,
	}

	expectedName := "Name Receita Federal"
	expectedBirthDay := "12/12/2012"

	client := &models.Client{
		Validated: true,
		Name:      expectedName,
		CPF:       input.CPF,
		Birthdate: input.BirthDate,
		Creators:  []*models.Creator{creator},
	}

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", "000.000.000-00", "12/12/2012").
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedName,
				NumeroDeCPF:    input.CPF,
				DataNascimento: expectedBirthDay,
			},
		}, nil)

	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).
		On("FindCreatorByUserEmail", creator.Contact.Email).
		Return(creator, nil)

	suite.mockClientRepository.(*MockClientRepository).
		On("FindByEmail", input.Email).
		Return(client, nil)

	_, err := suite.sut.CreateClient(input)

	suite.Error(err)
	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindCreatorByUserEmail", creator.Contact.Email)
	suite.mockClientRepository.(*MockClientRepository).AssertCalled(suite.T(), "FindByEmail", client.Contact.Email)
}

func TestClientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
