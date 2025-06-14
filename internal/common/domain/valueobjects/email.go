package common_vo

import (
	"fmt"
	"net/mail"
	"strings"
)

// Email represents an email address value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object. Returns error if the email is invalid.
func NewEmail(value string) (*Email, error) {
	// Trim spaces and convert to lowercase
	value = strings.TrimSpace(strings.ToLower(value))

	// Validate email format
	addr, err := mail.ParseAddress(value)
	if err != nil {
		return nil, fmt.Errorf("invalid email format: %w", err)
	}

	// Additional validation rules
	if len(addr.Address) > 254 { // RFC 5321
		return nil, fmt.Errorf("email address too long")
	}

	// Split email into local and domain parts
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid email format")
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part length (RFC 5321)
	if len(localPart) > 64 {
		return nil, fmt.Errorf("local part too long")
	}

	// Validate domain length (RFC 5321)
	if len(domain) > 255 {
		return nil, fmt.Errorf("domain too long")
	}

	return &Email{value: addr.Address}, nil
}

// Value returns the email value as a string
func (e *Email) Value() string {
	return e.value
}

// String returns the email address
func (e *Email) String() string {
	return e.value
}

// Equal checks if two emails are equal
func (e *Email) Equal(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

// Domain returns the domain part of the email address
func (e *Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email address
func (e *Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}
