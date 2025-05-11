package repositories

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Save(user *models.User) {
	database.DB.Create(&user)
	log.Default().Printf("Created new user with ID: %d, EMAIL: %s", user.ID, user.Email)
}

func (ur *UserRepository) FindByEmail(emailUser string) *models.User {
	var user *models.User
	database.DB.Preload("Creator").First(&user, "email = ?", emailUser)

	return user
}
