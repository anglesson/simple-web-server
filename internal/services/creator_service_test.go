package services_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreatorServiceTestSuite struct {
	suite.Suite
	sut             services.CreatorService
	mockCreatorRepo repositories.CreatorRepository
	mockRFService   gov.ReceitaFederalService
}

func TestCreatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CreatorServiceTestSuite))
}

func (suite *CreatorServiceTestSuite) SetupTest() {
	suite.mockCreatorRepo = new(MockCreatorRepository)
	suite.mockRFService = new(MockRFService)
	suite.sut = services.NewCreatorService(suite.mockCreatorRepo)
}

func (suite *CreatorServiceTestSuite) TestCreateCreator() {
	input := services.InputCreateCreator{
		Name:        "Valid Name",
		BirthDate:   "2012-12-12",
		PhoneNumber: "(12) 94567-8901",
		Email:       "valid@mail.com",
		CPF:         "058.997.950-77",
	}

	expectedCreator, _ := domain.NewCreator(
		input.Name,
		input.Email,
		input.CPF,
		input.PhoneNumber,
		input.BirthDate,
	)

	suite.mockCreatorRepo.(*MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(input)

	suite.NoError(err)
	suite.mockCreatorRepo.(*MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
}
