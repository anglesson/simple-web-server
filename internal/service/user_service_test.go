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
	inputUsername := "Any Username"
	inputPassword := "AnyPassword"
	inputEmail := "any@user.com"

	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", inputPassword).Return("HashedPassword")
	expectedUser, _ := domain.NewUser(inputUsername, inputEmail, "HashedPassword")

	suite.mockUserRepository.(*mocks.MockUserRepository).On("Save", &expectedUser).Return(nil)
	_, err := suite.sut.CreateUser(inputUsername, inputEmail, inputPassword)

	suite.NoError(err)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Save", &expectedUser)
}

func (suite *UserServiceTestSuite) TestCreateUserWithHashedPassword() {
	inputUsername := "Any Username"
	inputPassword := "AnyPassword"
	inputEmail := "any@user.com"

	expectedUser, _ := domain.NewUser(inputUsername, inputEmail, "HashedPassword")

	suite.mockEncrypter.(*utilsMocks.MockEncrypter).On("HashPassword", inputPassword).Return("HashedPassword")
	suite.mockUserRepository.(*mocks.MockUserRepository).On("Save", mock.Anything).Return(nil)
	user, err := suite.sut.CreateUser(inputUsername, inputEmail, inputPassword)

	suite.NoError(err)
	suite.mockEncrypter.(*utilsMocks.MockEncrypter).AssertCalled(suite.T(), "HashPassword", "AnyPassword")
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Save", &expectedUser)
	suite.Assert().Equal("HashedPassword", user.Password)
}
