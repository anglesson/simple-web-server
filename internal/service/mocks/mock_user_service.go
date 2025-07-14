package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(input service.InputCreateUser) (*domain.User, error) {
	args := m.Called(input)
	return args.Get(0).(*domain.User), args.Error(1)
}
