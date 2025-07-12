package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *domain.User) error {
	args := m.Called(&user)
	return args.Error(0)
}
func (m *MockUserRepository) FindByEmail(emailUser string) *domain.User {
	args := m.Called(emailUser)
	return args.Get(0).(*domain.User)
}
