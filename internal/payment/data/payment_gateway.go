package data

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/account"
)

type PaymentGateway interface {
	CreateAccount(name string, sellerID uint) (string, error)
}

type StripeService struct {
	secretKey string
}

func NewStripeService(secretKey string) *StripeService {
	return &StripeService{
		secretKey: secretKey,
	}
}

// CreateAccount creates a new Stripe account.
func (s *StripeService) CreateAccount(name string, sellerID uint) (string, error) {
	newAccount, err := account.New(&stripe.AccountParams{
		Controller: &stripe.AccountControllerParams{
			StripeDashboard: &stripe.AccountControllerStripeDashboardParams{
				Type: stripe.String("express"),
			},
			Fees: &stripe.AccountControllerFeesParams{
				Payer: stripe.String("application"),
			},
			Losses: &stripe.AccountControllerLossesParams{
				Payments: stripe.String("application"),
			},
		},
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			Name: stripe.String(name),
		},
	})

	if err != nil {
		fmt.Printf("An error occurred when calling the Stripe API to create an account: %v", err)
		return "", err
	}

	return newAccount.ID, nil
}
