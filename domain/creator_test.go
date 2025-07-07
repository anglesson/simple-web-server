package domain_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/domain"
)

func TestNewCreator(t *testing.T) {
	type InputType struct {
		Name      string
		Email     string
		CPF       string
		Phone     string
		Birthdate string
	}
	tests := []struct {
		name    string
		input   InputType
		wantErr bool
	}{
		{
			name: "Success",
			input: InputType{
				Name:      "Any Name",
				Email:     "valid@mail.com",
				CPF:       "058.997.950-77",
				Phone:     "(12) 98765-4321",
				Birthdate: "2004-01-01",
			},
			wantErr: false,
		},
		{
			name: "Invalid Name",
			input: InputType{
				Name:      "",
				Email:     "valid@mail.com",
				CPF:       "058.997.950-77",
				Phone:     "(12) 98765-4321",
				Birthdate: "2004-01-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid Email",
			input: InputType{
				Name:      "Valid Name",
				Email:     "invalid_mail.com",
				CPF:       "058.997.950-77",
				Phone:     "(12) 98765-4321",
				Birthdate: "2004-01-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid CPF",
			input: InputType{
				Name:      "Valid Name",
				Email:     "valid@mail.com",
				CPF:       "000.000.000-00",
				Phone:     "(12) 98765-4321",
				Birthdate: "2004-01-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid Phone",
			input: InputType{
				Name:      "Valid Name",
				Email:     "valid@mail.com",
				CPF:       "000.000.000-00",
				Phone:     "invalid phone",
				Birthdate: "2004-01-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid Birthdate",
			input: InputType{
				Name:      "Valid Name",
				Email:     "valid@mail.com",
				CPF:       "058.997.950-77",
				Phone:     "(12) 98765-4321",
				Birthdate: "invalid date",
			},
			wantErr: true,
		},
		{
			name: "Creator is minor",
			input: InputType{
				Name:      "Valid Name",
				Email:     "valid@mail.com",
				CPF:       "058.997.950-77",
				Phone:     "(12) 98765-4321",
				Birthdate: "2023-01-01",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator, err := domain.NewCreator(tt.input.Name, tt.input.Email, tt.input.CPF, tt.input.Phone, tt.input.Birthdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCreator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && creator.Name != tt.input.Name {
				t.Errorf("NewCreator() = %v, want %v", creator.Name, tt.input)
			}

			if !tt.wantErr && creator.Contact.Email.Value() != tt.input.Email {
				t.Errorf("NewCreator() = %v, want %v", creator.Contact.Email.Value(), tt.input.Email)
			}

			if !tt.wantErr && creator.CPF.String() != tt.input.CPF {
				t.Errorf("NewCreator() = %v, want %v", creator.CPF.String(), tt.input.CPF)
			}

			if !tt.wantErr && creator.Contact.Phone.String() != tt.input.Phone {
				t.Errorf("NewCreator() = %v, want %v", creator.Contact.Phone.String(), tt.input.Phone)
			}

			if !tt.wantErr && creator.Birthdate.String() != tt.input.Birthdate {
				t.Errorf("NewCreator() = %v, want %v", creator.Birthdate.String(), tt.input.Birthdate)
			}

			if !tt.wantErr && creator.Birthdate.IsAdult() == false {
				t.Errorf("NewCreator() = %v, want %v", creator.Birthdate.IsAdult(), true)
			}
		})
	}
}
