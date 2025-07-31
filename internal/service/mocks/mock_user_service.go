package mocks

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(input service.InputCreateUser) (*models.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(input service.InputLogin) (*models.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) RequestPasswordReset(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}

func (m *MockUserService) FindByEmail(email string) *models.User {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.User)
}
