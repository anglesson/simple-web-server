package service_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	mocks_repo "github.com/anglesson/simple-web-server/internal/repository/mocks"
	"github.com/anglesson/simple-web-server/internal/service"
	mocksService "github.com/anglesson/simple-web-server/internal/service/mocks"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/anglesson/simple-web-server/pkg/gov/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

var inputCreateCreator = service.InputCreateCreator{}

type CreatorServiceTestSuite struct {
	suite.Suite
	sut             service.CreatorService
	mockCreatorRepo repository.CreatorRepository
	mockRFService   gov.ReceitaFederalService
	mockUserService service.UserService
}

func TestCreatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CreatorServiceTestSuite))
}

func (suite *CreatorServiceTestSuite) SetupTest() {
	inputCreateCreator = service.InputCreateCreator{
		Name:                 "Valid Name",
		BirthDate:            "1990-12-12",
		PhoneNumber:          "(12) 94567-8901",
		Email:                "valid@mail.com",
		CPF:                  "058.997.950-77",
		Password:             "ValidPassword123!",
		PasswordConfirmation: "ValidPassword123!",
	}
	suite.mockCreatorRepo = new(mocks_repo.MockCreatorRepository)
	suite.mockRFService = new(mocks.MockRFService)
	suite.mockUserService = new(mocksService.MockUserService)
	suite.sut = service.NewCreatorService(suite.mockCreatorRepo, suite.mockRFService, suite.mockUserService)
}

func (suite *CreatorServiceTestSuite) TestCreateCreator() {
	expectedCreator, _ := domain.NewCreator(
		inputCreateCreator.Name,
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(nil, nil)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldCallUserService() {
	expectedCreator, _ := domain.NewCreator(
		inputCreateCreator.Name,
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(nil, nil)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockUserService.(*mocksService.MockUserService).
		AssertCalled(
			suite.T(),
			"CreateUser",
			inputCreateCreator.Name,
			inputCreateCreator.Email,
			inputCreateCreator.Password,
			inputCreateCreator.PasswordConfirmation,
		)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldUpdateCreatorWithReceitaFederalData() {
	expectedCreator, _ := domain.NewCreator(
		"Name RF",
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(nil, nil)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	_, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldThrowErrorIfCreatorHasARegister() {
	expectedCreator, _ := domain.NewCreator(
		inputCreateCreator.Name,
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(expectedCreator, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindCreatorByCPF", inputCreateCreator.CPF).
		Return(expectedCreator, nil)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	creator, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertNotCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldThrowErrorIfDataNotExistsInReceitaFederal() {
	expectedCreator, _ := domain.NewCreator(
		inputCreateCreator.Name,
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(expectedCreator, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindCreatorByCPF", inputCreateCreator.CPF).
		Return(expectedCreator, nil)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", expectedCreator.CPF.String(), expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: false, // Simulate data not found in Receita Federal
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	creator, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertNotCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldThrowErrorIfAnyDataIsInvalid() {
	inputCreateCreator.Email = "invalid_mail" // invalid mail

	expectedCreator, _ := domain.NewCreator(
		inputCreateCreator.Name,
		inputCreateCreator.Email,
		inputCreateCreator.CPF,
		inputCreateCreator.PhoneNumber,
		inputCreateCreator.BirthDate,
	)

	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF", "05899795077", "12/12/1990").
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       "Name RF",
				NumeroDeCPF:    "058.997.950-77",
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).
		On("CreateUser", inputCreateCreator.Name, inputCreateCreator.Email, inputCreateCreator.Password, inputCreateCreator.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("Save", expectedCreator).
		Return(nil)

	creator, err := suite.sut.CreateCreator(inputCreateCreator)

	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "FindByFilter")
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save")
	suite.mockRFService.(*mocks.MockRFService).AssertNotCalled(suite.T(), "ConsultaCPF")
}
