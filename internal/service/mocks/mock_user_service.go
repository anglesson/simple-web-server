package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(username, email, password, passwordConfirmation string) (*domain.User, error) {
	args := m.Called(username, email, password, passwordConfirmation)
	return args.Get(0).(*domain.User), args.Error(1)
}
