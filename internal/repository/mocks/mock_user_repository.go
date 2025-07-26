package mocks

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUserEmail(emailUser string) *models.User {
	args := m.Called(emailUser)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.User)
}
