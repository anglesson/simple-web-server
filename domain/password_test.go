package domain_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid password",
			input:   "Password123!",
			wantErr: false,
		},
		{
			name:    "password too short",
			input:   "Pass1!",
			wantErr: true,
		},
		{
			name:    "no uppercase",
			input:   "password123!",
			wantErr: true,
		},
		{
			name:    "no lowercase",
			input:   "PASSWORD123!",
			wantErr: true,
		},
		{
			name:    "no digit",
			input:   "Password!!",
			wantErr: true,
		},
		{
			name:    "no special character",
			input:   "Password123",
			wantErr: true,
		},
		{
			name:    "empty password",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewPassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
