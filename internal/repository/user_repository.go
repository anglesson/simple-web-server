package repository

import (
	"github.com/anglesson/simple-web-server/domain"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/database"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByUserEmail(emailUser string) *domain.User
}

type UserRepositoryImpl struct {
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Save(user *models.User) error {
	result := database.DB.Save(user)
	if result.Error != nil {
		log.Printf("Erro ao salvar usuário: %v", result.Error)
		return result.Error
	}

	log.Printf("Usuário atualizado com sucesso. ID: %d, EMAIL: %s", user.ID, user.Email)
	return nil
}

// TODO: add error handler
func (r *UserRepositoryImpl) FindByEmail(emailUser string) *models.User {
	var user models.User
	result := database.DB.Where("email = ?", emailUser).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por email: %v", result.Error)
		return nil
	}
	return &user
}

func (r *UserRepositoryImpl) FindBySessionToken(token string) *models.User {
	var user models.User
	result := database.DB.Where("session_token = ?", token).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por session token: %v", result.Error)
		return nil
	}
	return &user
}

func (r *UserRepositoryImpl) FindByStripeCustomerID(customerID string) *models.User {
	var user models.User
	err := database.DB.Where("stripe_customer_id = ?", customerID).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by Stripe customer ID: %v", err)
		return nil
	}
	return &user
}

func (r *UserRepositoryImpl) Create(user *domain.User) error {
	return nil
}

func (r *UserRepositoryImpl) FindByUserEmail(emailUser string) *domain.User {
	return nil
}
