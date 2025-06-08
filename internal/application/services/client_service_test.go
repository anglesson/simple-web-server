package services_test

import (
	"github.com/anglesson/simple-web-server/internal/application/dtos"
	"github.com/anglesson/simple-web-server/internal/application/ports"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var _ ports.ClientRepositoryPort = (*MockClientRepository)(nil)
var _ ports.CreatorRepositoryPort = (*MockCreatorRepository)(nil)

type MockCreatorRepository struct {
	mock.Mock
}

func (m *MockCreatorRepository) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) FindCreatorByUserEmail(email string) (*models.Creator, error) {
	args := m.Called()
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) Create(creator *models.Creator) error {
	args := m.Called()
	return args.Error(1)
}

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Save(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) FindClientsByCreator(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByIDAndCreators(client *models.Client, clientID uint, creator string) error {
	args := m.Called(client, clientID, creator)
	return args.Error(0)
}

func (m *MockClientRepository) FindByClientsWhereEbookNotSend(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error) {
	args := m.Called(creator, query)
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) FindByClientsWhereEbookWasSend(creator *models.Creator, query dtos.ClientQuery) (*[]models.Client, error) {
	args := m.Called()
	return args.Get(0).(*[]models.Client), args.Error(1)
}

func (m *MockClientRepository) InsertBatch(clients []*models.Client) error {
	args := m.Called()
	return args.Error(0)
}

type ClientServiceTestSuite struct {
	suite.Suite
	sut                   *services.ClientService
	mockClientRepository  ports.ClientRepositoryPort
	mockCreatorRepository ports.CreatorRepositoryPort
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockClientRepository = new(MockClientRepository)
	suite.mockCreatorRepository = new(MockCreatorRepository)
	suite.sut = services.NewClientService(suite.mockClientRepository, suite.mockCreatorRepository)
}
