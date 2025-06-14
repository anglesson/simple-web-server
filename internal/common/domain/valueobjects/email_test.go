package common_vo_test

import (
	"testing"

	common_vo "github.com/anglesson/simple-web-server/internal/common/domain/valueobjects"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid email",
			input:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with uppercase",
			input:   "USER@EXAMPLE.COM",
			wantErr: false,
		},
		{
			name:    "valid email with spaces",
			input:   " user@example.com ",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			input:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			input:   "user@",
			wantErr: true,
		},
		{
			name:    "invalid email - no local part",
			input:   "@example.com",
			wantErr: true,
		},
		{
			name:    "invalid email - empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - too long",
			input:   "a" + string(make([]byte, 255)) + "@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := common_vo.NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && email == nil {
				t.Error("NewEmail() returned nil Email when no error was expected")
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, err := common_vo.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("NewEmail() error = %v", err)
	}

	expected := "user@example.com"
	if got := email.String(); got != expected {
		t.Errorf("Email.String() = %v, want %v", got, expected)
	}
}

func TestEmail_Equal(t *testing.T) {
	email1, _ := common_vo.NewEmail("user@example.com")
	email2, _ := common_vo.NewEmail("user@example.com")
	email3, _ := common_vo.NewEmail("other@example.com")

	tests := []struct {
		name     string
		email1   *common_vo.Email
		email2   *common_vo.Email
		expected bool
	}{
		{
			name:     "equal emails",
			email1:   email1,
			email2:   email2,
			expected: true,
		},
		{
			name:     "different emails",
			email1:   email1,
			email2:   email3,
			expected: false,
		},
		{
			name:     "nil email",
			email1:   email1,
			email2:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.email1.Equal(tt.email2); got != tt.expected {
				t.Errorf("Email.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEmail_DomainAndLocalPart(t *testing.T) {
	email, err := common_vo.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("NewEmail() error = %v", err)
	}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "domain",
			got:      email.Domain(),
			expected: "example.com",
		},
		{
			name:     "local part",
			got:      email.LocalPart(),
			expected: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Email.%s() = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}
