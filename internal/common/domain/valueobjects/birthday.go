package common_vo

import (
	"fmt"
	"time"
)

// BirthDay represents a person's date of birth
type BirthDay struct {
	value time.Time
}

// NewBirthDay creates a new BirthDay value object. Returns error if the date is invalid.
func NewBirthDay(year, month, day int) (*BirthDay, error) {
	// Validate year
	currentYear := time.Now().Year()
	if year < 1900 || year > currentYear {
		return nil, fmt.Errorf("invalid year: must be between 1900 and %d", currentYear)
	}

	// Create time.Time object
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate if the date is valid (e.g., February 30th would be invalid)
	if date.Year() != year || int(date.Month()) != month || date.Day() != day {
		return nil, fmt.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	// Validate if the date is not in the future
	if date.After(time.Now()) {
		return nil, fmt.Errorf("birth date cannot be in the future")
	}

	return &BirthDay{value: date}, nil
}

// NewBirthDayFromTime creates a new BirthDay from a time.Time value
func NewBirthDayFromTime(t time.Time) (*BirthDay, error) {
	return NewBirthDay(t.Year(), int(t.Month()), t.Day())
}

// Value returns the birth date as time.Time
func (b *BirthDay) Value() time.Time {
	return b.value
}

// String returns the birth date in ISO format (YYYY-MM-DD)
func (b *BirthDay) String() string {
	return b.value.Format("2006-01-02")
}

// Equal checks if two birth dates are equal
func (b *BirthDay) Equal(other *BirthDay) bool {
	if other == nil {
		return false
	}
	return b.value.Equal(other.value)
}

// Age returns the current age in years
func (b *BirthDay) Age() int {
	now := time.Now()
	age := now.Year() - b.value.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < b.value.Month() || (now.Month() == b.value.Month() && now.Day() < b.value.Day()) {
		age--
	}

	return age
}

// IsAdult returns true if the person is 18 years or older
func (b *BirthDay) IsAdult() bool {
	return b.Age() >= 18
}

// Format returns the birth date formatted according to the given layout
func (b *BirthDay) Format(layout string) string {
	return b.value.Format(layout)
}

// Year returns the birth year
func (b *BirthDay) Year() int {
	return b.value.Year()
}

// Month returns the birth month (1-12)
func (b *BirthDay) Month() int {
	return int(b.value.Month())
}

// Day returns the birth day (1-31)
func (b *BirthDay) Day() int {
	return b.value.Day()
}
