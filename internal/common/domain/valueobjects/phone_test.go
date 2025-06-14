package common_vo

import (
	"testing"
)

func TestNewPhone(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid mobile number",
			input:   "(11) 98765-4321",
			wantErr: false,
		},
		{
			name:    "valid mobile number without formatting",
			input:   "11987654321",
			wantErr: false,
		},
		{
			name:    "invalid number - wrong length",
			input:   "1198765432",
			wantErr: true,
		},
		{
			name:    "invalid number - invalid area code",
			input:   "10987654321",
			wantErr: true,
		},
		{
			name:    "invalid number - not mobile (doesn't start with 9)",
			input:   "1187654321",
			wantErr: true,
		},
		{
			name:    "invalid number - all same digits",
			input:   "99999999999",
			wantErr: true,
		},
		{
			name:    "invalid number - empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid number - non-numeric characters",
			input:   "(11) ABC-1234",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && phone == nil {
				t.Error("NewPhone() returned nil Phone when no error was expected")
			}
		})
	}
}

func TestPhone_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "format valid number",
			input:    "11987654321",
			expected: "(11) 98765-4321",
		},
		{
			name:     "format already formatted number",
			input:    "(11) 98765-4321",
			expected: "(11) 98765-4321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.input)
			if err != nil {
				t.Fatalf("NewPhone() error = %v", err)
			}
			if got := phone.String(); got != tt.expected {
				t.Errorf("Phone.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPhone_Equal(t *testing.T) {
	phone1, _ := NewPhone("(11) 98765-4321")
	phone2, _ := NewPhone("(11) 98765-4321")
	phone3, _ := NewPhone("(21) 98765-4321")

	tests := []struct {
		name     string
		phone1   *Phone
		phone2   *Phone
		expected bool
	}{
		{
			name:     "equal phones",
			phone1:   phone1,
			phone2:   phone2,
			expected: true,
		},
		{
			name:     "different phones",
			phone1:   phone1,
			phone2:   phone3,
			expected: false,
		},
		{
			name:     "nil phone",
			phone1:   phone1,
			phone2:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.phone1.Equal(tt.phone2); got != tt.expected {
				t.Errorf("Phone.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPhone_AreaCodeAndNumber(t *testing.T) {
	phone, err := NewPhone("(11) 98765-4321")
	if err != nil {
		t.Fatalf("NewPhone() error = %v", err)
	}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "area code",
			got:      phone.AreaCode(),
			expected: "11",
		},
		{
			name:     "number",
			got:      phone.Number(),
			expected: "987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Phone.%s() = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestPhone_IsMobile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "mobile number",
			input:    "(11) 98765-4321",
			expected: true,
		},
		{
			name:     "landline number",
			input:    "(11) 3765-4321",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.input)
			if err != nil {
				// Skip test if the number is invalid (like landline)
				if tt.expected {
					t.Fatalf("NewPhone() error = %v", err)
				}
				return
			}
			if got := phone.IsMobile(); got != tt.expected {
				t.Errorf("Phone.IsMobile() = %v, want %v", got, tt.expected)
			}
		})
	}
}
