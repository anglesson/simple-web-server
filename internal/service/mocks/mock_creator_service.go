package mocks

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockCreatorService struct {
	mock.Mock
}

func (m MockCreatorService) CreateCreator(input service.InputCreateCreator) (*domain.Creator, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockCreatorService) FindCreatorByEmail(email string) (*models.Creator, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockCreatorService) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	//TODO implement me
	panic("implement me")
}
