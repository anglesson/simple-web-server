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

const (
	validName        = "Valid Name"
	validEmail       = "valid@mail.com"
	validCPF         = "058.997.950-77"
	validBirthDate   = "1990-12-12"
	validPhoneNumber = "(12) 94567-8901"
	validPassword    = "ValidPassword123!"
)

type CreatorServiceTestSuite struct {
	suite.Suite
	sut             service.CreatorService
	mockCreatorRepo repository.CreatorRepository
	mockRFService   gov.ReceitaFederalService
	mockUserService service.UserService
	testInput       service.InputCreateCreator
}

func TestCreatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CreatorServiceTestSuite))
}

func (suite *CreatorServiceTestSuite) SetupTest() {
	suite.setupTestInput()
	suite.setupMocks()
}

func (suite *CreatorServiceTestSuite) setupTestInput() {
	suite.testInput = service.InputCreateCreator{
		Name:                 validName,
		BirthDate:            validBirthDate,
		PhoneNumber:          validPhoneNumber,
		Email:                validEmail,
		CPF:                  validCPF,
		Password:             validPassword,
		PasswordConfirmation: validPassword,
	}
}

func (suite *CreatorServiceTestSuite) setupMocks() {
	suite.mockCreatorRepo = new(mocks_repo.MockCreatorRepository)
	suite.mockRFService = new(mocks.MockRFService)
	suite.mockUserService = new(mocksService.MockUserService)
	suite.sut = service.NewCreatorService(suite.mockCreatorRepo, suite.mockRFService, suite.mockUserService)
}

func (suite *CreatorServiceTestSuite) setupSuccessfulMockExpectations(expectedCreator *domain.Creator) {
	expectedFilter := domain.CreatorFilter{CPF: expectedCreator.CPF.Value()}

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedFilter).
		Return(nil, nil)
	suite.mockRFService.(*mocks.MockRFService).
		On("ConsultaCPF",
			expectedCreator.CPF.String(),
			expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: true,
			Result: gov.ConsultaData{
				NomeDaPF:       expectedCreator.Name,
				NumeroDeCPF:    validCPF,
				DataNascimento: "12/12/1990",
			},
		}, nil)

	suite.mockUserService.(*mocksService.MockUserService).On("CreateUser",
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.Password,
		suite.testInput.PasswordConfirmation).
		Return(&domain.User{}, nil)

	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).On("Save", expectedCreator).Return(nil)
}

func (suite *CreatorServiceTestSuite) TestCreateCreator_Success() {
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	suite.setupSuccessfulMockExpectations(expectedCreator)

	creator, err := suite.sut.CreateCreator(suite.testInput)

	suite.NoError(err)
	suite.NotNil(creator)
	suite.mockUserService.(*mocksService.MockUserService).AssertExpectations(suite.T())
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertExpectations(suite.T())
	suite.mockRFService.(*mocks.MockRFService).AssertExpectations(suite.T())
}

func (suite *CreatorServiceTestSuite) TestCreateCreator_ShouldCallUserService() {
	// Arrange
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	suite.setupSuccessfulMockExpectations(expectedCreator)

	// Act
	_, err := suite.sut.CreateCreator(suite.testInput)

	// Assert
	suite.NoError(err)
	suite.mockUserService.(*mocksService.MockUserService).
		AssertCalled(suite.T(), "CreateUser",
			suite.testInput.Name,
			suite.testInput.Email,
			suite.testInput.Password,
			suite.testInput.PasswordConfirmation)
}

func (suite *CreatorServiceTestSuite) TestCreateCreator_ShouldUpdateCreatorWithReceitaFederalData() {
	// Arrange
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	expectedCreator.Name = "Name Receita Federal"
	suite.setupSuccessfulMockExpectations(expectedCreator)

	// Act
	_, err := suite.sut.CreateCreator(suite.testInput)

	// Assert
	suite.NoError(err)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestCreateCreator_ShouldThrowErrorIfCreatorHasARegister() {
	// Arrange
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	expectedCreatorFilter := domain.CreatorFilter{
		CPF: expectedCreator.CPF.Value(),
	}
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).
		On("FindByFilter", expectedCreatorFilter).
		Return(expectedCreator, nil)
	suite.setupSuccessfulMockExpectations(expectedCreator)

	// Act
	creator, err := suite.sut.CreateCreator(suite.testInput)

	// Assert
	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertCalled(suite.T(), "FindByFilter", expectedCreatorFilter)
	suite.mockRFService.(*mocks.MockRFService).AssertNotCalled(suite.T(), "ConsultaCPF")
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save")
}

func (suite *CreatorServiceTestSuite) TestCreateCreator_FailsIfDataNotExistsInReceitaFederal() {
	// Arrange
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	suite.mockRFService.(*mocks.MockRFService).On("ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006")).
		Return(&gov.ReceitaFederalResponse{
			Status: false,
			Result: gov.ConsultaData{},
		}, nil)
	suite.setupSuccessfulMockExpectations(expectedCreator)

	// Act
	creator, err := suite.sut.CreateCreator(suite.testInput)

	// Assert
	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save", expectedCreator)
	suite.mockRFService.(*mocks.MockRFService).AssertCalled(
		suite.T(),
		"ConsultaCPF",
		expectedCreator.CPF.String(),
		expectedCreator.Birthdate.Format("02/01/2006"),
	)
}

func (suite *CreatorServiceTestSuite) TestShouldThrowErrorIfAnyDataIsInvalid() {
	// Arrange
	expectedCreator, _ := domain.NewCreator(
		suite.testInput.Name,
		suite.testInput.Email,
		suite.testInput.CPF,
		suite.testInput.PhoneNumber,
		suite.testInput.BirthDate,
	)
	suite.setupSuccessfulMockExpectations(expectedCreator)
	suite.testInput.Email = "invalid_mail" // invalid mail

	// Act
	creator, err := suite.sut.CreateCreator(suite.testInput)

	// Assert
	suite.Error(err)
	suite.Assert().Nil(creator)
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "FindByFilter")
	suite.mockCreatorRepo.(*mocks_repo.MockCreatorRepository).AssertNotCalled(suite.T(), "Save")
	suite.mockRFService.(*mocks.MockRFService).AssertNotCalled(suite.T(), "ConsultaCPF")
}
