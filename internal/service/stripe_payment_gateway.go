package service

import (
	"errors"
	"log"
)

type StripePaymentGateway struct {
	stripeService *StripeService
}

func NewStripePaymentGateway(stripeService *StripeService) PaymentGateway {
	return &StripePaymentGateway{
		stripeService: stripeService,
	}
}

func (spg *StripePaymentGateway) CreateCustomer(email, name string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}
	if name == "" {
		return "", errors.New("name is required")
	}

	customerID, err := spg.stripeService.CreateCustomer(email, name)
	if err != nil {
		log.Printf("Error creating customer in Stripe: %v", err)
		return "", err
	}

	return customerID, nil
}

func (spg *StripePaymentGateway) CreateSubscription(customerID, priceID string) (string, error) {
	if customerID == "" {
		return "", errors.New("customer ID is required")
	}
	if priceID == "" {
		return "", errors.New("price ID is required")
	}

	err := spg.stripeService.CreateSubscription(customerID, priceID)
	if err != nil {
		log.Printf("Error creating subscription in Stripe: %v", err)
		return "", err
	}

	// Note: The current StripeService.CreateSubscription doesn't return subscription ID
	// This would need to be updated to return the subscription ID
	return "", nil
}

func (spg *StripePaymentGateway) CancelSubscription(subscriptionID string) error {
	if subscriptionID == "" {
		return errors.New("subscription ID is required")
	}

	err := spg.stripeService.CancelSubscription(subscriptionID)
	if err != nil {
		log.Printf("Error canceling subscription in Stripe: %v", err)
		return err
	}

	return nil
}
