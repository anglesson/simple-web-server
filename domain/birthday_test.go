package domain_test

import (
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/domain"
)

func TestNewBirthDay(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()

	tests := []struct {
		name    string
		year    int
		month   int
		day     int
		wantErr bool
	}{
		{
			name:    "valid birthday",
			year:    1990,
			month:   1,
			day:     1,
			wantErr: false,
		},
		{
			name:    "invalid year - too old",
			year:    1899,
			month:   1,
			day:     1,
			wantErr: true,
		},
		{
			name:    "invalid year - future",
			year:    currentYear + 1,
			month:   1,
			day:     1,
			wantErr: true,
		},
		{
			name:    "invalid month - too high",
			year:    1990,
			month:   13,
			day:     1,
			wantErr: true,
		},
		{
			name:    "invalid month - too low",
			year:    1990,
			month:   0,
			day:     1,
			wantErr: true,
		},
		{
			name:    "invalid day - too high",
			year:    1990,
			month:   1,
			day:     32,
			wantErr: true,
		},
		{
			name:    "invalid day - too low",
			year:    1990,
			month:   1,
			day:     0,
			wantErr: true,
		},
		{
			name:    "invalid date - February 30th",
			year:    1990,
			month:   2,
			day:     30,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			birthday, err := domain.NewBirthDay(tt.year, tt.month, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBirthDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && birthday == nil {
				t.Error("NewBirthDay() returned nil BirthDay when no error was expected")
			}
		})
	}
}

func TestBirthDay_String(t *testing.T) {
	birthday, err := domain.NewBirthDay(1990, 1, 1)
	if err != nil {
		t.Fatalf("NewBirthDay() error = %v", err)
	}

	expected := "1990-01-01"
	if got := birthday.String(); got != expected {
		t.Errorf("BirthDay.String() = %v, want %v", got, expected)
	}
}

func TestBirthDay_Equal(t *testing.T) {
	birthday1, _ := domain.NewBirthDay(1990, 1, 1)
	birthday2, _ := domain.NewBirthDay(1990, 1, 1)
	birthday3, _ := domain.NewBirthDay(1991, 1, 1)

	tests := []struct {
		name      string
		birthday1 *domain.BirthDay
		birthday2 *domain.BirthDay
		expected  bool
	}{
		{
			name:      "equal birthdays",
			birthday1: birthday1,
			birthday2: birthday2,
			expected:  true,
		},
		{
			name:      "different birthdays",
			birthday1: birthday1,
			birthday2: birthday3,
			expected:  false,
		},
		{
			name:      "nil birthday",
			birthday1: birthday1,
			birthday2: nil,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.birthday1.Equal(tt.birthday2); got != tt.expected {
				t.Errorf("BirthDay.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBirthDay_Age(t *testing.T) {
	now := time.Now()
	year := now.Year() - 25 // 25 years ago

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected int
	}{
		{
			name:     "birthday this year",
			year:     year,
			month:    int(now.Month()),
			day:      now.Day(),
			expected: 25,
		},
		{
			name:     "birthday passed this year",
			year:     year,
			month:    int(now.Month()) - 1,
			day:      1,
			expected: 25,
		},
		{
			name:     "birthday not yet this year",
			year:     year,
			month:    int(now.Month()) + 1,
			day:      1,
			expected: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			birthday, err := domain.NewBirthDay(tt.year, tt.month, tt.day)
			if err != nil {
				t.Fatalf("NewBirthDay() error = %v", err)
			}
			if got := birthday.Age(); got != tt.expected {
				t.Errorf("BirthDay.Age() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBirthDay_IsAdult(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected bool
	}{
		{
			name:     "adult - 18 years old",
			year:     now.Year() - 18,
			month:    int(now.Month()),
			day:      now.Day(),
			expected: true,
		},
		{
			name:     "adult - 19 years old",
			year:     now.Year() - 19,
			month:    int(now.Month()),
			day:      now.Day(),
			expected: true,
		},
		{
			name:     "not adult - 17 years old",
			year:     now.Year() - 17,
			month:    int(now.Month()),
			day:      now.Day(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			birthday, err := domain.NewBirthDay(tt.year, tt.month, tt.day)
			if err != nil {
				t.Fatalf("NewBirthDay() error = %v", err)
			}
			if got := birthday.IsAdult(); got != tt.expected {
				t.Errorf("BirthDay.IsAdult() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBirthDay_Format(t *testing.T) {
	birthday, err := domain.NewBirthDay(1990, 1, 1)
	if err != nil {
		t.Fatalf("NewBirthDay() error = %v", err)
	}

	tests := []struct {
		name     string
		layout   string
		expected string
	}{
		{
			name:     "ISO format",
			layout:   "2006-01-02",
			expected: "1990-01-01",
		},
		{
			name:     "US format",
			layout:   "01/02/2006",
			expected: "01/01/1990",
		},
		{
			name:     "European format",
			layout:   "02/01/2006",
			expected: "01/01/1990",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := birthday.Format(tt.layout); got != tt.expected {
				t.Errorf("BirthDay.Format() = %v, want %v", got, tt.expected)
			}
		})
	}
}
