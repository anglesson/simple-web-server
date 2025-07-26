package gorm

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"gorm.io/gorm"
)

type SubscriptionGormRepository struct {
}

func NewSubscriptionGormRepository() repository.SubscriptionRepository {
	return &SubscriptionGormRepository{}
}

func (sr *SubscriptionGormRepository) Create(subscription *models.Subscription) error {
	err := database.DB.Create(subscription).Error
	if err != nil {
		log.Printf("Erro ao criar subscription: %s", err)
		return errors.New("erro ao criar subscription")
	}
	return nil
}

func (sr *SubscriptionGormRepository) FindByUserID(userID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := database.DB.Where("user_id = ?", userID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Erro ao buscar subscription por user_id: %s", err)
		return nil, errors.New("erro ao buscar subscription")
	}
	return &subscription, nil
}

func (sr *SubscriptionGormRepository) FindByStripeCustomerID(customerID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := database.DB.Where("stripe_customer_id = ?", customerID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Erro ao buscar subscription por stripe_customer_id: %s", err)
		return nil, errors.New("erro ao buscar subscription")
	}
	return &subscription, nil
}

func (sr *SubscriptionGormRepository) FindByStripeSubscriptionID(subscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := database.DB.Where("stripe_subscription_id = ?", subscriptionID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Erro ao buscar subscription por stripe_subscription_id: %s", err)
		return nil, errors.New("erro ao buscar subscription")
	}
	return &subscription, nil
}

func (sr *SubscriptionGormRepository) Update(subscription *models.Subscription) error {
	err := database.DB.Model(subscription).Updates(subscription).Error
	if err != nil {
		log.Printf("Erro ao atualizar subscription: %s", err)
		return errors.New("erro ao atualizar subscription")
	}
	return nil
}

func (sr *SubscriptionGormRepository) Save(subscription *models.Subscription) error {
	err := database.DB.Save(subscription).Error
	if err != nil {
		log.Printf("Erro ao salvar subscription: %s", err)
		return errors.New("erro ao salvar subscription")
	}
	return nil
}
