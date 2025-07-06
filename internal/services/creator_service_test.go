package services_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repositories"
	mocks_repo "github.com/anglesson/simple-web-server/internal/repositories/mocks"
	"github.com/anglesson/simple-web-server/internal/services"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/anglesson/simple-web-server/pkg/gov/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreatorServiceTestSuite struct {
	suite.Suite
	sut             services.CreatorService
	mockCreatorRepo repositories.CreatorRepository
	mockRFService   gov.ReceitaFederalService
}

func (suite *CreatorServiceTestSuite) SetupTest() {
	suite.mockCreatorRepo = new(mocks_repo.MockCreatorRepository)
	suite.mockRFService = new(mocks.MockRFService)
	suite.sut = services.NewCreatorService(suite.mockCreatorRepo, suite.mockRFService)
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

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/2012",
			},
		}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(input)

	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
}

func (suite *CreatorServiceTestSuite) TestShouldUpdateCreatorWithReceitaFederalData() {
	input := services.InputCreateCreator{
		Name:        "Valid Name",
		BirthDate:   "2012-12-12",
		PhoneNumber: "(12) 94567-8901",
		Email:       "valid@mail.com",
		CPF:         "058.997.950-77",
	}

	expectedCreator, _ := domain.NewCreator(
		"Name RF",
		input.Email,
		input.CPF,
		input.PhoneNumber,
		input.BirthDate,
	)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/2012",
			},
		}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(input)

	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func TestCreatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CreatorServiceTestSuite))
}
