package domain_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			name     string
			cpf      string
			birthDay string
			email    string
			phone    string
		}
		wantErr bool
	}{
		{
			name: "valid client",
			input: struct {
				name     string
				cpf      string
				birthDay string
				email    string
				phone    string
			}{
				name:     "John Doe",
				cpf:      "461.371.640-39",
				birthDay: "1990-01-01",
				email:    "john@example.com",
				phone:    "1234567890",
			},
			wantErr: false,
		},
		{
			name: "invalid CPF",
			input: struct {
				name     string
				cpf      string
				birthDay string
				email    string
				phone    string
			}{
				name:     "John Doe",
				cpf:      "123.456.789-00",
				birthDay: "1990-01-01",
				email:    "john@example.com",
				phone:    "1234567890",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewClient(tt.input.name, tt.input.cpf, tt.input.birthDay, tt.input.email, tt.input.phone)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Update(t *testing.T) {
	client, err := domain.NewClient("John Doe", "461.371.640-39", "1990-01-01", "john@example.com", "1234567890")
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	err = client.Update("Jane Doe", "461.371.640-39", "jane@example.com", "0987654321")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", client.Name)
	assert.Equal(t, "jane@example.com", client.Email)
	assert.Equal(t, "0987654321", client.Phone)
}
