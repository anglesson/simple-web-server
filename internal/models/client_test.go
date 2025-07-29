package models_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/models"
)

func TestClient_GetInitials(t *testing.T) {
	tests := []struct {
		name     string
		client   *models.Client
		expected string
	}{
		{
			name:     "Single name",
			client:   &models.Client{Name: "João"},
			expected: "J",
		},
		{
			name:     "Two names",
			client:   &models.Client{Name: "João Silva"},
			expected: "JS",
		},
		{
			name:     "Three names",
			client:   &models.Client{Name: "João Pedro Silva"},
			expected: "JS",
		},
		{
			name:     "Empty name",
			client:   &models.Client{Name: ""},
			expected: "?",
		},
		{
			name:     "Multiple spaces",
			client:   &models.Client{Name: "João   Silva"},
			expected: "JS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.client.GetInitials()
			if result != tt.expected {
				t.Errorf("GetInitials() = %v, want %v", result, tt.expected)
			}
		})
	}
}
