package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username             string     `json:"username" validate:"required"`
	Password             string     `json:"password" validate:"required"`
	Email                string     `json:"email" validate:"required,email" gorm:"unique"`
	TrialStartDate       *time.Time `json:"trial_start_date"`
	TrialEndDate         *time.Time `json:"trial_end_date"`
	IsTrialActive        bool       `json:"is_trial_active"`
	StripeCustomerID     string     `json:"stripe_customer_id"`
	StripeSubscriptionID string     `json:"stripe_subscription_id"`
	SubscriptionStatus   string     `json:"subscription_status"`
	SubscriptionEndDate  *time.Time `json:"subscription_end_date"`
	CSRFToken            string
	SessionToken         string
}

func NewUser(username, password, email string) *User {
	now := time.Now()
	trialEndDate := now.AddDate(0, 0, 7) // 7 days trial

	return &User{
		Username:       username,
		Password:       password,
		Email:          email,
		TrialStartDate: &now,
		TrialEndDate:   &trialEndDate,
		IsTrialActive:  true,
	}
}

func (u *User) GetInitials() string {
	words := strings.Fields(u.Username)
	if len(words) == 0 {
		return ""
	}
	initials := strings.ToUpper(string(words[0][0]))
	if len(words) > 1 {
		lastWord := words[len(words)-1]
		if len(lastWord) > 0 {
			initials += strings.ToUpper(string(lastWord[0]))
		}
	}
	return initials
}

func (u *User) IsInTrialPeriod() bool {
	if u.TrialEndDate == nil {
		return false
	}
	return u.IsTrialActive && time.Now().Before(*u.TrialEndDate)
}

func (u *User) DaysLeftInTrial() int {
	if u.TrialEndDate == nil {
		return 0
	}
	days := u.TrialEndDate.Sub(time.Now()).Hours() / 24
	if days < 0 {
		return 0
	}
	return int(days)
}

func (u *User) IsSubscribed() bool {
	if u.SubscriptionStatus == "active" {
		return true
	}
	if u.SubscriptionEndDate != nil && time.Now().Before(*u.SubscriptionEndDate) {
		return true
	}
	return false
}
