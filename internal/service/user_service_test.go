package service_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/mocks"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/utils"
	utilsMocks "github.com/anglesson/simple-web-server/pkg/utils/mocks"
	"github.com/stretchr/testify/mock"
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
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserRepository = new(mocks.MockUserRepository)
	suite.mockEncrypter = new(utilsMocks.MockEncrypter)
	suite.sut = service.NewUserService(suite.mockUserRepository, suite.mockEncrypter)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreateUser() {

	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByEmail", input.Email).Return(nil)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", input.Password).Return("HashedPassword123!")
	expectedUser, _ := domain.NewUser(input.Username, input.Email, "HashedPassword123!")

	suite.mockUserRepository.(*mocks.MockUserRepository).On("Save", &expectedUser).Return(nil)
	_, err := suite.sut.CreateUser(input.Username, input.Email, input.Password, input.PasswordConfirmation)

	suite.NoError(err)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Save", &expectedUser)
}

func (suite *UserServiceTestSuite) TestCreateUserWithHashedPassword() {
	expectedUser, _ := domain.NewUser(input.Username, input.Email, "HashedPassword123!")

	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByEmail", input.Email).Return(nil)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", input.Password).Return("HashedPassword123!")
	suite.mockUserRepository.(*mocks.MockUserRepository).On("Save", mock.Anything).Return(nil)
	user, err := suite.sut.CreateUser(input.Username, input.Email, input.Password, input.PasswordConfirmation)

	suite.NoError(err)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).AssertCalled(suite.T(), "HashPassword", input.Password)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Save", &expectedUser)
	suite.Assert().Equal("HashedPassword123!", user.Password.Value())
}

func (suite *UserServiceTestSuite) TestShouldReturnErrorIfPasswordAndConfirmationAreDifferent() {
	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByEmail", input.Email).Return(nil)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", input.Password).Return("HashedPassword")
	user, err := suite.sut.CreateUser(input.Username, input.Email, input.Password, input.PasswordConfirmation)

	suite.Error(err)
	suite.Nil(user)
}

func (suite *UserServiceTestSuite) TestShouldReturnErrorIfUserAlreadyExists() {
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", input.Password).Return("HashedPassword")
	suite.mockUserRepository.(*mocks.MockUserRepository).On("FindByEmail", input.Email).Return(&domain.User{})
	user, err := suite.sut.CreateUser(input.Username, input.Email, input.Password, input.PasswordConfirmation)

	suite.Error(err)
	suite.Nil(user)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "FindByEmail", input.Email)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertNotCalled(suite.T(), "Save")
}
