package domain_test

import (
	"github.com/anglesson/simple-web-server/domain"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	type InputType struct {
		Username string
		Email    string
		Password string
	}
	tests := []struct {
		name    string
		input   InputType
		wantErr bool
	}{
		{
			name: "Success",
			input: InputType{
				Username: "Any Username",
				Email:    "valid@mail.com",
				Password: "ValidPassword123!",
			},
			wantErr: false,
		}, {
			name: "invalid format email",
			input: InputType{
				Username: "Any Username",
				Email:    "invalid_valid_mail",
				Password: "ValidPassword123!",
			},
			wantErr: true,
		}, {
			name: "empty username",
			input: InputType{
				Username: "",
				Email:    "valid@mail.com",
				Password: "ValidPassword123!",
			},
			wantErr: true,
		}, {
			name: "empty email",
			input: InputType{
				Username: "Any username",
				Email:    "",
				Password: "ValidPassword123!",
			},
			wantErr: true,
		}, {
			name: "empty password",
			input: InputType{
				Username: "Any username",
				Email:    "valid@mail.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "max length 50 characters",
			input: InputType{
				Username: strings.Repeat("a", 51),
				Email:    "valid@mail.com",
				Password: "ValidPassword123!",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.input.Username, tt.input.Email, tt.input.Password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(tt.input.Username) > 50 && tt.input.Username[:50] != user.Username {
				t.Errorf("NewUser() = lenght %v, want lenthg %v", len(user.Username), len(tt.input.Username))
			}

			if !tt.wantErr && user.Username != tt.input.Username && len(tt.input.Username) < 50 {
				t.Errorf("NewUser() = %v, want %v", user.Username, tt.input.Username)
			}

			if !tt.wantErr && user.Email.Value() != tt.input.Email {
				t.Errorf("NewUser() = %v, want %v", user.Email.Value(), tt.input.Username)
			}

			if !tt.wantErr && !user.Password.Equals(tt.input.Password) {
				t.Errorf("NewUser() = %v, want %v", user.Password, tt.input.Password)
			}
		})
	}
}
