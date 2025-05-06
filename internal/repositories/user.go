package repositories

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

func Save(user *models.User) {
	database.DB.Create(&user)
	log.Default().Printf("Created new user with ID: %d, EMAIL: %s", user.ID, user.Email)
}

func FindByEmail(emailUser string) *models.User {
	var user *models.User
	database.DB.First(&user, "email = ?", emailUser)

	return user
}
