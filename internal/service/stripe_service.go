package service

import (
	"log"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
)

type StripeService struct {
	apiKey string
}

func NewStripeService() *StripeService {
	stripe.Key = config.AppConfig.StripeSecretKey
	return &StripeService{
		apiKey: config.AppConfig.StripeSecretKey,
	}
}

func (s *StripeService) CreateCustomer(user *models.User) error {
	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Username),
		Metadata: map[string]string{
			"user_id": strconv.FormatUint(uint64(user.ID), 10),
		},
	}

	c, err := customer.New(params)
	if err != nil {
		log.Printf("Error creating Stripe customer: %v", err)
		return err
	}

	user.StripeCustomerID = c.ID
	return nil
}

func (s *StripeService) CreateSubscription(customerID string, priceID string) error {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}

	_, err := subscription.New(params)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return err
	}

	return nil
}

func (s *StripeService) CancelSubscription(subscriptionID string) error {
	_, err := subscription.Cancel(subscriptionID, nil)
	if err != nil {
		log.Printf("Error canceling subscription: %v", err)
		return err
	}

	return nil
}
