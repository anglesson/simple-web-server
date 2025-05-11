package services

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

type CreatorService struct{}

func NewCreatorService() *CreatorService {
	return &CreatorService{}
}

func (cs *CreatorService) FindCreatorByEmail(email string) (*models.Creator, error) {
	userRepository := repositories.NewUserRepository()
	user := userRepository.FindByEmail(email)
	creator := models.Creator{
		UserID: user.ID,
	}
	result := database.DB.First(&creator)

	if result.Error != nil {
		log.Printf("Erro ao buscar creator: %s", result.Error)
		return nil, errors.New("creator not found")
	}

	return &creator, nil
}
