package mocks

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/service"
)

type MockCreatorService struct{}

func (m MockCreatorService) CreateCreator(input service.InputCreateCreator) (*models.Creator, error) {
	// Mock implementation
	return nil, nil
}

func (m MockCreatorService) FindCreatorByEmail(email string) (*models.Creator, error) {
	// Mock implementation
	return nil, nil
}

func (m MockCreatorService) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	// Mock implementation
	return nil, nil
}
