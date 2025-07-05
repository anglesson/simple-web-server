package gorm

import (
	"errors"
	"github.com/anglesson/simple-web-server/domain"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/database"
)

type CreatorRepository struct {
}

func (cr *CreatorRepository) Save(creator *domain.Creator) error {
	//TODO implement me
	panic("implement me")
}

func NewCreatorRepository() *CreatorRepository {
	return &CreatorRepository{}
}

func (cr *CreatorRepository) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	var creator models.Creator
	err := database.DB.
		First(&creator, "user_id = ?", userID).Error
	if err != nil {
		log.Printf("creator isn't recovery. error: %s", err.Error())
		return nil, errors.New("creator not found")
	}
	return &creator, nil
}

func (cr *CreatorRepository) FindCreatorByUserEmail(email string) (*models.Creator, error) {
	var creator models.Creator
	err := database.DB.
		Joins("JOIN users ON users.id = creators.user_id").
		First(&creator, "users.email = ?", email).Error
	if err != nil {
		log.Printf("creator isn't recovery. error: %s", err.Error())
		return nil, errors.New("creator not found")
	}
	return &creator, nil
}

func (cr *CreatorRepository) Create(creator *models.Creator) error {
	err := database.DB.Create(&creator).Error
	if err != nil {
		log.Printf("fail on create 'creator': %s", err.Error())
		return errors.New("creator not found")
	}
	return nil
}
