package service_test

import (
	mocks_repo "github.com/anglesson/simple-web-server/internal/repository/mocks"
	"github.com/anglesson/simple-web-server/pkg/gov/mocks"
	"testing"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/suite"
)

var _ repository.ClientRepository = (*mocks_repo.MockClientRepository)(nil)
var _ repository.CreatorRepository = (*mocks_repo.MockCreatorRepository)(nil)
var _ gov.ReceitaFederalService = (*mocks.MockRFService)(nil)

type ClientServiceTestSuite struct {
	suite.Suite
	sut                   service.ClientService
	mockClientRepository  repository.ClientRepository
	mockCreatorRepository repository.CreatorRepository
	mockRFService         gov.ReceitaFederalService
}

func (suite *ClientServiceTestSuite) SetupTest() {
	suite.mockClientRepository = new(mocks_repo.MockClientRepository)
	suite.mockCreatorRepository = new(mocks_repo.MockCreatorRepository)
	suite.mockRFService = new(mocks.MockRFService)
	suite.sut = service.NewClientService(suite.mockClientRepository, suite.mockCreatorRepository, suite.mockRFService)
}

func (suite *ClientServiceTestSuite) TestCreateClient() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}

	input := service.CreateClientInput{
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

	suite.mockClientRepository.(*mocks_repo.MockClientRepository).
		On("FindByEmail", input.Email).
		Return(nil, nil)

	suite.mockClientRepository.(*mocks_repo.MockClientRepository).
		On("Save", client).
		Return(nil)

	_, err := suite.sut.CreateClient(input)

	suite.NoError(err)
	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindCreatorByUserEmail", creator.Contact.Email)
	suite.mockClientRepository.(*mocks_repo.MockClientRepository).AssertCalled(suite.T(), "Save", client)
}

func (suite *ClientServiceTestSuite) TestShouldReturnErrorIfClientExists() {
	creator := &models.Creator{Contact: models.Contact{Email: "creator@mail.com"}}

	input := service.CreateClientInput{
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

	suite.mockClientRepository.(*mocks_repo.MockClientRepository).
		On("FindByEmail", input.Email).
		Return(client, nil)

	_, err := suite.sut.CreateClient(input)

	suite.Error(err)
	suite.mockCreatorRepository.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindCreatorByUserEmail", creator.Contact.Email)
	suite.mockClientRepository.(*mocks_repo.MockClientRepository).AssertCalled(suite.T(), "FindByEmail", client.Contact.Email)
}

func TestClientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ClientServiceTestSuite))
}
