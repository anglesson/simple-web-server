package domain

import (
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidID                 = errors.New("invalid ID")
	ErrInvalidUserID             = errors.New("invalid user ID")
	ErrInvalidPlanID             = errors.New("invalid plan ID")
	ErrInvalidTrialStartDate     = errors.New("invalid trial start date")
	ErrInvalidTrialEndDate       = errors.New("invalid trial end date")
	ErrInvalidCustomerID         = errors.New("invalid customer ID")
	ErrInvalidSubscriptionID     = errors.New("invalid subscription ID")
	ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")
	ErrInvalidOrigin             = errors.New("invalid origin")
	ErrInvalidCreatedAt          = errors.New("invalid created at timestamp")
	ErrInvalidUpdatedAt          = errors.New("invalid updated at timestamp")
)

type Subscription struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	PlanID              string    `json:"plan_id"`
	TrialStartDate      time.Time `json:"trial_start_date"`
	TrialEndDate        time.Time `json:"trial_end_date"`
	IsTrialActive       bool      `json:"is_trial_active"`
	CustomerID          string    `json:"stripe_customer_id"`
	SubscriptionID      string    `json:"stripe_subscription_id"`
	SubscriptionStatus  string    `json:"subscription_status"`
	SubscriptionEndDate time.Time `json:"subscription_end_date"`
	Origin              string    `json:"origin"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewSubscription(
	id string,
	userID string,
	planID string,
	trialStartDate time.Time,
	trialEndDate time.Time,
	isTrialActive bool,
	customerID string,
	subscriptionID string,
	subscriptionStatus string,
	subscriptionEndDate time.Time,
	origin string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Subscription, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if planID == "" {
		return nil, ErrInvalidPlanID
	}
	if trialStartDate.IsZero() {
		return nil, ErrInvalidTrialStartDate
	}
	if trialEndDate.IsZero() {
		return nil, ErrInvalidTrialEndDate
	}
	if customerID == "" {
		return nil, ErrInvalidCustomerID
	}
	if subscriptionID == "" {
		return nil, ErrInvalidSubscriptionID
	}
	if subscriptionStatus == "" {
		return nil, ErrInvalidSubscriptionStatus
	}
	if origin == "" {
		return nil, ErrInvalidOrigin
	}
	if createdAt.IsZero() {
		return nil, ErrInvalidCreatedAt
	}
	if updatedAt.IsZero() {
		return nil, ErrInvalidUpdatedAt
	}

	if trialEndDate.Before(trialStartDate) {
		return nil, ErrInvalidTrialEndDate
	}

	return &Subscription{
		ID:                  id,
		UserID:              userID,
		PlanID:              planID,
		TrialStartDate:      trialStartDate,
		TrialEndDate:        trialEndDate,
		IsTrialActive:       isTrialActive,
		CustomerID:          customerID,
		SubscriptionID:      subscriptionID,
		SubscriptionStatus:  subscriptionStatus,
		SubscriptionEndDate: subscriptionEndDate,
		Origin:              origin,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}, nil
}

func (s *Subscription) IsInTrialPeriod() bool {
	if s.TrialEndDate.IsZero() {
		return false
	}
	return s.IsTrialActive && time.Now().Before(s.TrialEndDate)
}

func (s *Subscription) DaysLeftInTrial() int {
	if s.TrialEndDate.IsZero() {
		return 0
	}
	days := s.TrialEndDate.Sub(time.Now()).Hours() / 24
	if math.Ceil(days) < 0 {
		return 0
	}
	return int(days)
}

func (s *Subscription) IsSubscribed() bool {
	if s.SubscriptionStatus == "active" {
		return true
	}
	if !s.SubscriptionEndDate.IsZero() && time.Now().Before(s.SubscriptionEndDate) {
		return true
	}
	return false
}

func (s *Subscription) CancelSubscription() {
	s.SubscriptionStatus = "canceled"
	s.UpdatedAt = time.Now()
}

func (s *Subscription) EndTrial() {
	s.IsTrialActive = false
	s.UpdatedAt = time.Now()
}
