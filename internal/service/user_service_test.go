package service_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/mocks"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/utils"
	utilsMocks "github.com/anglesson/simple-web-server/pkg/utils/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InputCreateUser struct {
	Username             string
	Password             string
	PasswordConfirmation string
	Email                string
}

var input = InputCreateUser{
	Username:             "Any Username",
	Password:             "Password123!",
	PasswordConfirmation: "Password123!",
	Email:                "any@user.com",
}

type UserServiceTestSuite struct {
	suite.Suite
	sut                service.UserService
	mockUserRepository repository.UserRepository
	mockEncrypter      utils.Encrypter
	testInput          InputCreateUser
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.setUpInput()
	suite.setupMocks()
}

func (suite *UserServiceTestSuite) setUpInput() {
	suite.testInput = InputCreateUser{
		Username:             "Valid UserName",
		Email:                "valid@mail.com",
		Password:             "Password123!",
		PasswordConfirmation: "Password123!",
	}
}

func (suite *UserServiceTestSuite) setupMocks() {
	suite.mockUserRepository = new(mocks.MockUserRepository)
	suite.mockEncrypter = new(utilsMocks.MockEncrypter)
	suite.sut = service.NewUserService(suite.mockUserRepository, suite.mockEncrypter)
}

func (suite *UserServiceTestSuite) setupSuccessfulMockExpectations(expectedUser *domain.User) {
	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByUserEmail", suite.testInput.Email).Return(nil)
	suite.mockUserRepository.(*mocks.MockUserRepository).On("Create", expectedUser).Return(nil)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", suite.testInput.Password).Return("HashedPassword123!")
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
	// Arrange
	expectedUser, _ := domain.NewUser(suite.testInput.Username, suite.testInput.Email, "HashedPassword123!")
	suite.setupSuccessfulMockExpectations(expectedUser)

	// Act
	user, err := suite.sut.CreateUser(
		suite.testInput.Username,
		suite.testInput.Email,
		suite.testInput.Password,
		suite.testInput.PasswordConfirmation)

	// Assert
	suite.NoError(err)
	suite.NotNil(user)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Create", expectedUser)
}

func (suite *UserServiceTestSuite) TestCreateUser_ShouldCallHashPassword() {
	// Arrange
	expectedUser, _ := domain.NewUser(suite.testInput.Username, suite.testInput.Email, "HashedPassword123!")
	suite.setupSuccessfulMockExpectations(expectedUser)

	// Act
	user, err := suite.sut.CreateUser(suite.testInput.Username, suite.testInput.Email, suite.testInput.Password, suite.testInput.PasswordConfirmation)

	// Assert
	suite.NoError(err)
	suite.NotNil(user)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).AssertCalled(suite.T(), "HashPassword", input.Password)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Create", expectedUser)
	suite.Assert().Equal("HashedPassword123!", user.Password.Value())
}

func (suite *UserServiceTestSuite) TestShouldReturnErrorIfPasswordAndConfirmationAreDifferent() {
	// Arrange
	expectedUser, _ := domain.NewUser(suite.testInput.Username, suite.testInput.Email, "HashedPassword123!")
	suite.testInput.PasswordConfirmation = "DifferentPassword"
	suite.setupSuccessfulMockExpectations(expectedUser)

	// Act
	user, err := suite.sut.CreateUser(suite.testInput.Username, suite.testInput.Email, suite.testInput.Password, suite.testInput.PasswordConfirmation)

	// Assert
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UserServiceTestSuite) TestShouldReturnErrorIfUserAlreadyExists() {
	// Arrange
	suite.setupSuccessfulMockExpectations(&domain.User{})
	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByUserEmail", input.Email).Return(&domain.User{})

	// Act
	user, err := suite.sut.CreateUser(input.Username, input.Email, input.Password, input.PasswordConfirmation)

	// Assert
	suite.Error(err)
	suite.Nil(user)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "FindByUserEmail", input.Email)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertNotCalled(suite.T(), "Create")
}
