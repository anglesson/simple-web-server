package repository

import (
	"log"

	"github.com/anglesson/simple-web-server/domain"
	"gorm.io/gorm"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/database"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByUserEmail(emailUser string) *domain.User
}

type GormUserRepositoryImpl struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepositoryImpl {
	return &GormUserRepositoryImpl{
		db: db,
	}
}

func (r *GormUserRepositoryImpl) Save(user *models.User) error {
	result := database.DB.Save(user)
	if result.Error != nil {
		log.Printf("Erro ao salvar usuário: %v", result.Error)
		return result.Error
	}

	log.Printf("Usuário atualizado com sucesso. ID: %d, EMAIL: %s", user.ID, user.Email)
	return nil
}

// TODO: add error handler
func (r *GormUserRepositoryImpl) FindByEmail(emailUser string) *models.User {
	var user models.User
	result := database.DB.Where("email = ?", emailUser).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por email: %v", result.Error)
		return nil
	}
	return &user
}

func (r *GormUserRepositoryImpl) FindBySessionToken(token string) *models.User {
	var user models.User
	result := database.DB.Where("session_token = ?", token).First(&user)
	if result.Error != nil {
		log.Printf("Erro ao buscar usuário por session token: %v", result.Error)
		return nil
	}
	return &user
}

func (r *GormUserRepositoryImpl) FindByStripeCustomerID(customerID string) *models.User {
	var user models.User
	err := database.DB.Where("stripe_customer_id = ?", customerID).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by Stripe customer ID: %v", err)
		return nil
	}
	return &user
}

func (r *GormUserRepositoryImpl) Create(user *domain.User) error {
	var userModel models.User
	userModel.Username = user.Username
	userModel.Email = user.Email.Value()
	userModel.Password = user.Password.Value()

	err := r.db.Create(&userModel).Error
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	user.ID = userModel.ID
	return nil
}

func (r *GormUserRepositoryImpl) FindByUserEmail(emailUser string) *domain.User {
	var user models.User
	err := r.db.Where("email = ?", emailUser).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil
	}
	userDomain, err := domain.NewUser(user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Error creating user domain: %v", err)
		return nil
	}
	userDomain.ID = user.ID
	return userDomain
}
