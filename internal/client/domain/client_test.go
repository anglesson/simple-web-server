package client_domain

import (
	"errors"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		cpf           string
		birthDay      string
		email         string
		phone         string
		expected      *Client
		expectedError error
	}{
		{
			name:     "John Doe",
			cpf:      "123.456.789-00",
			birthDay: "1990-01-01",
			email:    "john@example.com",
			phone:    "+55 11 99999-9999",
			expected: &Client{Name: "John Doe", CPF: "123.456.789-00", BirthDay: "1990-01-01", Email: "john@example.com", Phone: "+55 11 99999-9999"},
		},
		{
			name:     "Jane Smith",
			cpf:      "987.654.321-00",
			birthDay: "1990-01-01",
			email:    "jane@example.com",
			phone:    "+55 11 98888-8888",
			expected: &Client{Name: "Jane Smith", CPF: "987.654.321-00", BirthDay: "1990-01-01", Email: "jane@example.com", Phone: "+55 11 98888-8888"},
		},
		{
			name:          "",
			cpf:           "987.654.321-00",
			birthDay:      "1990-01-01",
			email:         "jane@example.com",
			phone:         "+55 11 98888-8888",
			expectedError: errors.New("name is required"),
		},
		{
			name:          "John Doe",
			cpf:           "", // empty CPF
			birthDay:      "1990-01-01",
			email:         "john@example.com",
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("CPF is required"),
		},
		{
			name:          "John Doe",
			cpf:           "123.456.789-00",
			birthDay:      "", // empty birthDay
			email:         "john@example.com",
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("birthday is required"),
		},
		{
			name:          "John Doe",
			cpf:           "123.456.789-00",
			birthDay:      "1990-01-01",
			email:         "", // empty email
			phone:         "+55 11 99999-9999",
			expectedError: errors.New("email is required"),
		},
		{
			name:          "John Doe",
			cpf:           "123.456.789-00",
			birthDay:      "1990-01-01",
			email:         "john@example.com",
			phone:         "", // empty phone
			expectedError: errors.New("phone is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.name, tt.cpf, tt.birthDay, tt.email, tt.phone)

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

			if client.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", client.Name, tt.expected.Name)
			}
			if client.CPF != tt.expected.CPF {
				t.Errorf("CPF = %v, want %v", client.CPF, tt.expected.CPF)
			}
			if client.BirthDay != tt.expected.BirthDay {
				t.Errorf("BirthDay = %v, want %v", client.BirthDay, tt.expected.BirthDay)
			}
			if client.Email != tt.expected.Email {
				t.Errorf("Email = %v, want %v", client.Email, tt.expected.Email)
			}
			if client.Phone != tt.expected.Phone {
				t.Errorf("Phone = %v, want %v", client.Phone, tt.expected.Phone)
			}
		})
	}
}
