package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockCreatorRepository struct {
	mock.Mock
}

func (m *MockCreatorRepository) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) FindCreatorByUserEmail(email string) (*models.Creator, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorRepository) FindByFilter(query domain.CreatorFilter) (*domain.Creator, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*domain.Creator), args.Error(1)
}

func (m *MockCreatorRepository) Create(creator *models.Creator) error {
	args := m.Called(creator)
	return args.Error(0)
}

func (m *MockCreatorRepository) Save(creator *domain.Creator) error {
	args := m.Called(creator)
	return args.Error(0)
}
