package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string     `json:"username" validate:"required"`
	Password         string     `json:"password" validate:"required"`
	Email            string     `json:"email" validate:"required,email" gorm:"unique"`
	TrialStartDate   *time.Time `json:"trial_start_date"`
	TrialEndDate     *time.Time `json:"trial_end_date"`
	IsTrialActive    bool       `json:"is_trial_active"`
	StripeCustomerID string     `json:"stripe_customer_id"`
	CSRFToken        string
	SessionToken     string
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
	initials := ""
	for _, word := range words {
		if len(word) > 0 {
			initials += strings.ToUpper(string(word[0]))
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
