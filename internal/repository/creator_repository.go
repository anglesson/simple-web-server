package repository

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/models"
)

type CreatorRepository interface {
	FindCreatorByUserID(userID uint) (*models.Creator, error)
	FindCreatorByUserEmail(email string) (*models.Creator, error)
	FindByFilter(query domain.CreatorFilter) (*domain.Creator, error)
	Create(creator *models.Creator) error
	Save(creator *domain.Creator) error
}
