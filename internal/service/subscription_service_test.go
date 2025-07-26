package service

import (
	"errors"
	"testing"

	"github.com/anglesson/simple-web-server/internal/models"
	repoMocks "github.com/anglesson/simple-web-server/internal/repository/mocks"
	govMocks "github.com/anglesson/simple-web-server/pkg/gov/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionService_CreateSubscription(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		planID         string
		setupMocks     func(*repoMocks.MockSubscriptionRepository)
		expectedError  bool
		expectedResult *models.Subscription
	}{
		{
			name:   "should create subscription successfully",
			userID: 1,
			planID: "test_plan",
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.Subscription")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:           "should return error when user ID is zero",
			userID:         0,
			planID:         "test_plan",
			setupMocks:     func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:           "should return error when plan ID is empty",
			userID:         1,
			planID:         "",
			setupMocks:     func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			planID: "test_plan",
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.Subscription")).Return(errors.New("database error"))
			},
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repoMocks.MockSubscriptionRepository)
			mockRF := new(govMocks.MockRFService)
			tt.setupMocks(mockRepo)

			service := NewSubscriptionService(mockRepo, mockRF)

			result, err := service.CreateSubscription(tt.userID, tt.planID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.UserID)
				assert.Equal(t, tt.planID, result.PlanID)
				assert.True(t, result.IsTrialActive)
				assert.Equal(t, "web", result.Origin)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionService_FindByUserID(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		setupMocks     func(*repoMocks.MockSubscriptionRepository)
		expectedError  bool
		expectedResult *models.Subscription
	}{
		{
			name:   "should find subscription successfully",
			userID: 1,
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				subscription := &models.Subscription{
					UserID: 1,
					PlanID: "test_plan",
				}
				mockRepo.On("FindByUserID", uint(1)).Return(subscription, nil)
			},
			expectedError: false,
		},
		{
			name:           "should return error when user ID is zero",
			userID:         0,
			setupMocks:     func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("FindByUserID", uint(1)).Return(nil, errors.New("database error"))
			},
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repoMocks.MockSubscriptionRepository)
			mockRF := new(govMocks.MockRFService)
			tt.setupMocks(mockRepo)

			service := NewSubscriptionService(mockRepo, mockRF)

			result, err := service.FindByUserID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionService_ActivateSubscription(t *testing.T) {
	tests := []struct {
		name                 string
		subscription         *models.Subscription
		stripeCustomerID     string
		stripeSubscriptionID string
		setupMocks           func(*repoMocks.MockSubscriptionRepository)
		expectedError        bool
	}{
		{
			name: "should activate subscription successfully",
			subscription: &models.Subscription{
				UserID: 1,
				PlanID: "test_plan",
			},
			stripeCustomerID:     "cus_123",
			stripeSubscriptionID: "sub_123",
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("Save", mock.AnythingOfType("*models.Subscription")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:                 "should return error when subscription is nil",
			subscription:         nil,
			stripeCustomerID:     "cus_123",
			stripeSubscriptionID: "sub_123",
			setupMocks:           func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:        true,
		},
		{
			name: "should return error when customer ID is empty",
			subscription: &models.Subscription{
				UserID: 1,
				PlanID: "test_plan",
			},
			stripeCustomerID:     "",
			stripeSubscriptionID: "sub_123",
			setupMocks:           func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:        true,
		},
		{
			name: "should return error when subscription ID is empty",
			subscription: &models.Subscription{
				UserID: 1,
				PlanID: "test_plan",
			},
			stripeCustomerID:     "cus_123",
			stripeSubscriptionID: "",
			setupMocks:           func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repoMocks.MockSubscriptionRepository)
			mockRF := new(govMocks.MockRFService)
			tt.setupMocks(mockRepo)

			service := NewSubscriptionService(mockRepo, mockRF)

			err := service.ActivateSubscription(tt.subscription, tt.stripeCustomerID, tt.stripeSubscriptionID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.stripeCustomerID, tt.subscription.StripeCustomerID)
				assert.Equal(t, tt.stripeSubscriptionID, tt.subscription.StripeSubscriptionID)
				assert.Equal(t, "active", tt.subscription.SubscriptionStatus)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionService_CancelSubscription(t *testing.T) {
	tests := []struct {
		name          string
		subscription  *models.Subscription
		setupMocks    func(*repoMocks.MockSubscriptionRepository)
		expectedError bool
	}{
		{
			name: "should cancel subscription successfully",
			subscription: &models.Subscription{
				UserID: 1,
				PlanID: "test_plan",
			},
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("Save", mock.AnythingOfType("*models.Subscription")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "should return error when subscription is nil",
			subscription:  nil,
			setupMocks:    func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repoMocks.MockSubscriptionRepository)
			mockRF := new(govMocks.MockRFService)
			tt.setupMocks(mockRepo)

			service := NewSubscriptionService(mockRepo, mockRF)

			err := service.CancelSubscription(tt.subscription)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "canceled", tt.subscription.SubscriptionStatus)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionService_EndTrial(t *testing.T) {
	tests := []struct {
		name          string
		subscription  *models.Subscription
		setupMocks    func(*repoMocks.MockSubscriptionRepository)
		expectedError bool
	}{
		{
			name: "should end trial successfully",
			subscription: &models.Subscription{
				UserID:        1,
				PlanID:        "test_plan",
				IsTrialActive: true,
			},
			setupMocks: func(mockRepo *repoMocks.MockSubscriptionRepository) {
				mockRepo.On("Save", mock.AnythingOfType("*models.Subscription")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "should return error when subscription is nil",
			subscription:  nil,
			setupMocks:    func(mockRepo *repoMocks.MockSubscriptionRepository) {},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repoMocks.MockSubscriptionRepository)
			mockRF := new(govMocks.MockRFService)
			tt.setupMocks(mockRepo)

			service := NewSubscriptionService(mockRepo, mockRF)

			err := service.EndTrial(tt.subscription)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.False(t, tt.subscription.IsTrialActive)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
