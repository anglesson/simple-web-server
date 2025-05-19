package services

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
)

type CreatorService struct {
	creatorRepository *repositories.CreatorRepository
}

func NewCreatorService() *CreatorService {
	return &CreatorService{
		creatorRepository: repositories.NewCreatorRepository(),
	}
}

func (cs *CreatorService) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserID(userID)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}

func (cs *CreatorService) FindCreatorByEmail(email string) (*models.Creator, error) {
	creator, err := cs.creatorRepository.FindCreatorByUserEmail(email)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}
