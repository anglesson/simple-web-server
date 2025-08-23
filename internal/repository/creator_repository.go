package repository

import (
	"github.com/anglesson/simple-web-server/internal/models"
)

type CreatorRepository interface {
	FindCreatorByUserID(userID string) (*models.Creator, error)
	FindByCPF(cpf string) (*models.Creator, error)
	FindByID(id uint) (*models.Creator, error)
	Create(creator *models.Creator) error
}
