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

func (ur *UserRepository) Save(user *models.User) error {
	result := database.DB.Save(user)
	if result.Error != nil {
		log.Printf("Erro ao salvar usuário: %v", result.Error)
		return result.Error
	}

	log.Printf("Usuário atualizado com sucesso. ID: %d, EMAIL: %s", user.ID, user.Email)
	return nil
}

// TODO: add error handler
func (ur *UserRepository) FindByEmail(emailUser string) *models.User {
	var user models.User
	result := database.DB.Where("email = ?", emailUser).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por email: %v", result.Error)
		return nil
	}
	return &user
}

func (r *UserRepository) FindBySessionToken(token string) *models.User {
	var user models.User
	result := database.DB.Where("session_token = ?", token).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por session token: %v", result.Error)
		return nil
	}
	return &user
}
