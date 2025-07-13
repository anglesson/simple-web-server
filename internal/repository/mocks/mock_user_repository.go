package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(&user)
	return args.Error(0)
}
func (m *MockUserRepository) FindByUserEmail(emailUser string) *domain.User {
	args := m.Called(emailUser)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.User)
}
