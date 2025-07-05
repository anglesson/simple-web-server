package services

import (
	"errors"
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
)

type CreatorService interface {
	CreateCreator(input InputCreateCreator) (*domain.Creator, error)
	FindCreatorByEmail(email string) (*models.Creator, error)
	FindCreatorByUserID(userID uint) (*models.Creator, error)
}

type InputCreateCreator struct {
	Name        string `json:"name"`
	CPF         string `json:"cpf"`
	BirthDate   string `json:"birthDate"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type CreatorServiceImpl struct {
	creatorRepo repositories.CreatorRepository
}

func NewCreatorService(creatorRepo repositories.CreatorRepository) CreatorService {
	return &CreatorServiceImpl{
		creatorRepo: creatorRepo,
	}
}

func (cs *CreatorServiceImpl) CreateCreator(input InputCreateCreator) (*domain.Creator, error) {
	creator, err := domain.NewCreator(
		input.Name,
		input.Email,
		input.CPF,
		input.PhoneNumber,
		input.BirthDate,
	)
	if err != nil {
		return nil, err
	}

	err = cs.creatorRepo.Save(creator)
	if err != nil {
		return nil, err
	}
	return creator, nil
}

func (cs *CreatorServiceImpl) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	creator, err := cs.creatorRepo.FindCreatorByUserID(userID)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}

func (cs *CreatorServiceImpl) FindCreatorByEmail(email string) (*models.Creator, error) {
	creator, err := cs.creatorRepo.FindCreatorByUserEmail(email)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}
