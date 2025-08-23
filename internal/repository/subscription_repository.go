package repository

import (
	"github.com/anglesson/simple-web-server/internal/models"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	FindByUserID(userID string) (*models.Subscription, error)
	FindByStripeCustomerID(customerID string) (*models.Subscription, error)
	FindByStripeSubscriptionID(subscriptionID string) (*models.Subscription, error)
	Update(subscription *models.Subscription) error
	Save(subscription *models.Subscription) error
}
