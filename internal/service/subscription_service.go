package service

import (
	"errors"
	"time"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/gov"
)

type SubscriptionService interface {
	CreateSubscription(userID uint, planID string) (*models.Subscription, error)
	FindByUserID(userID uint) (*models.Subscription, error)
	FindByStripeCustomerID(customerID string) (*models.Subscription, error)
	FindByStripeSubscriptionID(subscriptionID string) (*models.Subscription, error)
	ActivateSubscription(subscription *models.Subscription, stripeCustomerID, stripeSubscriptionID string) error
	UpdateSubscriptionStatus(subscription *models.Subscription, status string, endDate *time.Time) error
	CancelSubscription(subscription *models.Subscription) error
	EndTrial(subscription *models.Subscription) error
}

type subscriptionServiceImpl struct {
	subscriptionRepository repository.SubscriptionRepository
	receitaFederalService  gov.ReceitaFederalService
}

func NewSubscriptionService(
	subscriptionRepository repository.SubscriptionRepository,
	receitaFederalService gov.ReceitaFederalService,
) SubscriptionService {
	return &subscriptionServiceImpl{
		subscriptionRepository: subscriptionRepository,
		receitaFederalService:  receitaFederalService,
	}
}

func (ss *subscriptionServiceImpl) CreateSubscription(userID uint, planID string) (*models.Subscription, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}
	if planID == "" {
		return nil, errors.New("plan ID is required")
	}

	subscription := models.NewSubscription(userID, planID)

	err := ss.subscriptionRepository.Create(subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (ss *subscriptionServiceImpl) FindByUserID(userID uint) (*models.Subscription, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	return ss.subscriptionRepository.FindByUserID(userID)
}

func (ss *subscriptionServiceImpl) FindByStripeCustomerID(customerID string) (*models.Subscription, error) {
	if customerID == "" {
		return nil, errors.New("customer ID is required")
	}

	return ss.subscriptionRepository.FindByStripeCustomerID(customerID)
}

func (ss *subscriptionServiceImpl) FindByStripeSubscriptionID(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, errors.New("subscription ID is required")
	}

	return ss.subscriptionRepository.FindByStripeSubscriptionID(subscriptionID)
}

func (ss *subscriptionServiceImpl) ActivateSubscription(subscription *models.Subscription, stripeCustomerID, stripeSubscriptionID string) error {
	if subscription == nil {
		return errors.New("subscription is required")
	}
	if stripeCustomerID == "" {
		return errors.New("stripe customer ID is required")
	}
	if stripeSubscriptionID == "" {
		return errors.New("stripe subscription ID is required")
	}

	subscription.ActivateSubscription(stripeCustomerID, stripeSubscriptionID)

	return ss.subscriptionRepository.Save(subscription)
}

func (ss *subscriptionServiceImpl) UpdateSubscriptionStatus(subscription *models.Subscription, status string, endDate *time.Time) error {
	if subscription == nil {
		return errors.New("subscription is required")
	}
	if status == "" {
		return errors.New("status is required")
	}

	subscription.UpdateSubscriptionStatus(status, endDate)

	return ss.subscriptionRepository.Save(subscription)
}

func (ss *subscriptionServiceImpl) CancelSubscription(subscription *models.Subscription) error {
	if subscription == nil {
		return errors.New("subscription is required")
	}

	subscription.CancelSubscription()

	return ss.subscriptionRepository.Save(subscription)
}

func (ss *subscriptionServiceImpl) EndTrial(subscription *models.Subscription) error {
	if subscription == nil {
		return errors.New("subscription is required")
	}

	subscription.EndTrial()

	return ss.subscriptionRepository.Save(subscription)
}
