package domain_test

import (
	"errors"
	"testing"

	"github.com/anglesson/simple-web-server/internal/domain"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		cpf           string
		birthDay      string
		email         string
		phone         string
		expected      *map[string]string
		expectedError error
	}{
		{
			name:          "Invalid CPF format",
			cpf:           "123.456.789-XX", // invalid format
			birthDay:      "1990-01-01",
			email:         "john@example.com",
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("CPF inválido para o cliente: CPF inválido (dígitos verificadores)"),
		},
		{
			name:          "Empty CPF",
			cpf:           "", // empty CPF
			birthDay:      "1990-01-01",
			email:         "john@example.com",
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("CPF inválido para o cliente: CPF deve ter 11 dígitos"),
		},
		{
			name:     "John Doe",
			cpf:      " 201.476.380-11 ", // valid CPF with extra spaces
			birthDay: "1990-01-01",
			email:    "john@example.com",
			phone:    "+55 11 99999-9999",
			expected: &map[string]string{"Name": "John Doe", "CPF": "201.476.380-11", "BirthDay": "1990-01-01", "Email": "john@example.com", "Phone": "+55 11 99999-9999"},
		},
		{
			name:          "Invalid CPF checksum",
			cpf:           "123.456.789-11", // invalid checksum
			birthDay:      "1990-01-01",
			email:         "john@example.com",
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("CPF inválido para o cliente: CPF inválido (dígitos verificadores)"),
		},
		{
			name:     "John Doe",
			cpf:      "235.489.640-95", // valid CPF
			birthDay: "1990-01-01",
			email:    "john@example.com",
			phone:    "+55 11 99999-9999",
			expected: &map[string]string{"Name": "John Doe", "CPF": "235.489.640-95", "BirthDay": "1990-01-01", "Email": "john@example.com", "Phone": "+55 11 99999-9999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := domain.NewClient(tt.name, tt.cpf, tt.birthDay, tt.email, tt.phone)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					if client != nil {
						t.Error("Client should be nil")
					}
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if client.Name != (*tt.expected)["Name"] {
				t.Errorf("Name = %v, want %v", client.Name, (*tt.expected)["Name"])
			}
			if client.CPF.String() != (*tt.expected)["CPF"] {
				t.Errorf("CPF = %v, want %v", client.CPF.String(), (*tt.expected)["CPF"])
			}
			if client.BirthDay.Value().Format("2006-01-02") != (*tt.expected)["BirthDay"] {
				t.Errorf("BirthDay = %v, want %v", client.BirthDay, (*tt.expected)["BirthDay"])
			}
			if client.Email != (*tt.expected)["Email"] {
				t.Errorf("Email = %v, want %v", client.Email, (*tt.expected)["Email"])
			}
			if client.Phone != (*tt.expected)["Phone"] {
				t.Errorf("Phone = %v, want %v", client.Phone, (*tt.expected)["Phone"])
			}
		})
	}
}
