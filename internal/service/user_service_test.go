package service_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/mocks"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	sut                service.UserService
	mockUserRepository repository.UserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserRepository = new(mocks.MockUserRepository)
	suite.sut = service.NewUserService(suite.mockUserRepository)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	inputUsername := "Any Username"
	inputPassword := "AnyPassword"
	inputEmail := "any@user.com"

	expectedUser, _ := domain.NewUser(inputUsername, inputEmail, inputPassword)

	suite.mockUserRepository.(*mocks.MockUserRepository).On("Save", &expectedUser).Return(nil)
	_, err := suite.sut.CreateUser(inputUsername, inputEmail, inputPassword)

	suite.NoError(err)
	suite.mockUserRepository.(*mocks.MockUserRepository).AssertCalled(suite.T(), "Save", &expectedUser)
}
