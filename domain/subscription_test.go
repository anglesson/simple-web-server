package domain_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/domain"
)

func TestNewSubscription(t *testing.T) {
	type InputSubscriptionType struct {
		ID                  string
		UserID              string
		PlanID              string
		TrialStartDate      int64
		TrialEndDate        int64
		IsTrialActive       bool
		CustomerID          string
		SubscriptionID      string
		SubscriptionStatus  string
		SubscriptionEndDate int64
		Origin              string
		CreatedAt           int64
		UpdatedAt           int64
	}

	tests := []struct {
		name    string
		input   InputSubscriptionType
		want    *domain.Subscription
		wantErr bool
	}{
		{
			name: "Success",
			input: InputSubscriptionType{
				ID:                  "sub_123",
				UserID:              "user_456",
				PlanID:              "plan_789",
				TrialStartDate:      1633046400,
				TrialEndDate:        1635724800,
				IsTrialActive:       true,
				CustomerID:          "customer_001",
				SubscriptionID:      "sub_123",
				SubscriptionStatus:  "active",
				SubscriptionEndDate: 1638316800,
				Origin:              "web",
				CreatedAt:           1633046400,
				UpdatedAt:           1633046400,
			},
			want: &domain.Subscription{
				ID:                  "sub_123",
				UserID:              "user_456",
				PlanID:              "plan_789",
				TrialStartDate:      time.Unix(1633046400, 0),
				TrialEndDate:        time.Unix(1635724800, 0),
				IsTrialActive:       true,
				CustomerID:          "customer_001",
				SubscriptionID:      "sub_123",
				SubscriptionStatus:  "active",
				SubscriptionEndDate: time.Unix(1638316800, 0),
				Origin:              "web",
				CreatedAt:           time.Unix(1633046400, 0),
				UpdatedAt:           time.Unix(1633046400, 0),
			},
			wantErr: false,
		},
		{
			name: "Error when missing required fields",
			input: InputSubscriptionType{
				ID:        "",
				UserID:    "",
				PlanID:    "",
				Origin:    "web",
				CreatedAt: 1633046400,
				UpdatedAt: 1633046400,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error when TrialEndDate is before TrialStartDate",
			input: InputSubscriptionType{
				ID:                  "sub_124",
				UserID:              "user_457",
				PlanID:              "plan_790",
				TrialStartDate:      1635724800,
				TrialEndDate:        1633046400,
				IsTrialActive:       true,
				CustomerID:          "customer_002",
				SubscriptionID:      "sub_124",
				SubscriptionStatus:  "active",
				SubscriptionEndDate: 1638316800,
				Origin:              "web",
				CreatedAt:           1633046400,
				UpdatedAt:           1633046400,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewSubscription(
				tt.input.ID,
				tt.input.UserID,
				tt.input.PlanID,
				time.Unix(tt.input.TrialStartDate, 0),
				time.Unix(tt.input.TrialEndDate, 0),
				tt.input.IsTrialActive,
				tt.input.CustomerID,
				tt.input.SubscriptionID,
				tt.input.SubscriptionStatus,
				time.Unix(tt.input.SubscriptionEndDate, 0),
				tt.input.Origin,
				time.Unix(tt.input.CreatedAt, 0),
				time.Unix(tt.input.UpdatedAt, 0),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriptionMethods(t *testing.T) {
	now := time.Now()
	oneDay := 24 * time.Hour

	sub := &domain.Subscription{
		ID:                  "sub_001",
		UserID:              "user_001",
		PlanID:              "plan_001",
		TrialStartDate:      now.Add(oneDay),
		TrialEndDate:        now.Add(oneDay * 2),
		IsTrialActive:       true,
		CustomerID:          "cust_001",
		SubscriptionID:      "sub_001",
		SubscriptionStatus:  "active",
		SubscriptionEndDate: now.Add(10 * oneDay),
		Origin:              "web",
		CreatedAt:           now.Add(-2 * oneDay),
		UpdatedAt:           now.Add(-2 * oneDay),
	}

	t.Run("IsInTrialPeriod true", func(t *testing.T) {
		if !sub.IsInTrialPeriod() {
			t.Errorf("IsInTrialPeriod() = false, want true")
		}
	})

	t.Run("IsInTrialPeriod false when trial ended", func(t *testing.T) {
		originalEndDate := sub.TrialEndDate
		sub.TrialEndDate = now.Add(-oneDay)

		t.Cleanup(func() {
			sub.TrialEndDate = originalEndDate // Esta linha será executada SEMPRE ao final deste sub-teste
		})

		if sub.IsInTrialPeriod() {
			t.Errorf("IsInTrialPeriod() = true, want false")
		}
		sub.TrialEndDate = now.Add(oneDay) // restore
	})

	t.Run("DaysLeftInTrial positive", func(t *testing.T) {
		days := sub.DaysLeftInTrial()
		if days < 1 {
			t.Errorf("DaysLeftInTrial() = %d, want >= 1", days)
		}
	})

	t.Run("DaysLeftInTrial zero when trial ended", func(t *testing.T) {
		sub.TrialEndDate = now.Add(-oneDay)
		if sub.DaysLeftInTrial() != 0 {
			t.Errorf("DaysLeftInTrial() = %d, want 0", sub.DaysLeftInTrial())
		}
		sub.TrialEndDate = now.Add(oneDay) // restore
	})

	t.Run("IsSubscribed true when active", func(t *testing.T) {
		sub.SubscriptionStatus = "active"
		if !sub.IsSubscribed() {
			t.Errorf("IsSubscribed() = false, want true")
		}
	})

	t.Run("IsSubscribed true when end date in future", func(t *testing.T) {
		sub.SubscriptionStatus = "canceled"
		sub.SubscriptionEndDate = now.Add(oneDay)
		if !sub.IsSubscribed() {
			t.Errorf("IsSubscribed() = false, want true")
		}
	})

	t.Run("IsSubscribed false when expired", func(t *testing.T) {
		sub.SubscriptionStatus = "canceled"
		sub.SubscriptionEndDate = now.Add(-oneDay)
		if sub.IsSubscribed() {
			t.Errorf("IsSubscribed() = true, want false")
		}
		sub.SubscriptionEndDate = now.Add(10 * oneDay) // restore
	})

	t.Run("CancelSubscription sets status and updatedAt", func(t *testing.T) {
		sub.SubscriptionStatus = "active"
		oldUpdated := sub.UpdatedAt
		sub.CancelSubscription()
		if sub.SubscriptionStatus != "canceled" {
			t.Errorf("CancelSubscription() did not set status to canceled")
		}
		if !sub.UpdatedAt.After(oldUpdated) {
			t.Errorf("CancelSubscription() did not update UpdatedAt")
		}
	})

	t.Run("EndTrial sets IsTrialActive and updatedAt", func(t *testing.T) {
		sub.IsTrialActive = true
		oldUpdated := sub.UpdatedAt
		sub.EndTrial()
		if sub.IsTrialActive {
			t.Errorf("EndTrial() did not set IsTrialActive to false")
		}
		if !sub.UpdatedAt.After(oldUpdated) {
			t.Errorf("EndTrial() did not update UpdatedAt")
		}
	})
}
