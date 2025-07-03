package domain

import (
	"fmt"
	"regexp"
)

// Phone represents a Brazilian phone number value object
type Phone struct {
	value string
}

// NewPhone creates a new Phone value object. Returns error if the phone number is invalid.
func NewPhone(value string) (*Phone, error) {
	// Clean the phone number (remove all non-digit characters)
	cleanNumber := cleanPhone(value)

	// Validate the phone number format
	if err := validatePhone(cleanNumber); err != nil {
		return nil, err
	}

	return &Phone{value: cleanNumber}, nil
}

// Value returns the phone number value as a string (only digits)
func (p *Phone) Value() string {
	return p.value
}

// String returns the formatted phone number (XX) XXXXX-XXXX
func (p *Phone) String() string {
	if len(p.value) != 11 {
		return p.value
	}
	return fmt.Sprintf("(%s) %s-%s",
		p.value[0:2],
		p.value[2:7],
		p.value[7:11],
	)
}

// Equal checks if two phone numbers are equal
func (p *Phone) Equal(other *Phone) bool {
	if other == nil {
		return false
	}
	return p.value == other.value
}

// AreaCode returns the area code (DDD) of the phone number
func (p *Phone) AreaCode() string {
	if len(p.value) < 2 {
		return ""
	}
	return p.value[0:2]
}

// Number returns the phone number without area code
func (p *Phone) Number() string {
	if len(p.value) < 11 {
		return p.value
	}
	return p.value[2:]
}

// IsMobile returns true if the phone number is a mobile number
func (p *Phone) IsMobile() bool {
	if len(p.value) < 11 {
		return false
	}
	// In Brazil, mobile numbers start with 9 after the area code
	return p.value[2] == '9'
}

// cleanPhone removes all non-digit characters from the phone number
func cleanPhone(phone string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(phone, "")
}

// validatePhone validates the phone number format
func validatePhone(phone string) error {
	// Check if the phone number has the correct length
	if len(phone) != 11 {
		return fmt.Errorf("invalid phone number length: must have 11 digits")
	}

	// Validate area code (DDD)
	areaCode := phone[0:2]
	if !isValidAreaCode(areaCode) {
		return fmt.Errorf("invalid area code: %s", areaCode)
	}

	// Validate if it's a mobile number (must start with 9)
	if phone[2] != '9' {
		return fmt.Errorf("invalid phone number: mobile numbers must start with 9")
	}

	// Validate if all digits are the same (invalid number)
	if allPhoneDigitsSame(phone) {
		return fmt.Errorf("invalid phone number: all digits are the same")
	}

	return nil
}

// isValidAreaCode checks if the area code is valid in Brazil
func isValidAreaCode(areaCode string) bool {
	// List of valid Brazilian area codes
	validAreaCodes := map[string]bool{
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"21": true, "22": true, "24": true, "27": true, "28": true,
		"31": true, "32": true, "33": true, "34": true, "35": true, "37": true, "38": true,
		"41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"51": true, "53": true, "54": true, "55": true,
		"61": true, "62": true, "63": true, "64": true, "65": true, "66": true, "67": true, "68": true, "69": true,
		"71": true, "73": true, "74": true, "75": true, "77": true, "79": true,
		"81": true, "82": true, "83": true, "84": true, "85": true, "86": true, "87": true, "88": true, "89": true,
		"91": true, "92": true, "93": true, "94": true, "95": true, "96": true, "97": true, "98": true, "99": true,
	}

	return validAreaCodes[areaCode]
}

// allPhoneDigitsSame checks if all digits in the phone number are the same
func allPhoneDigitsSame(phone string) bool {
	first := phone[0]
	for i := 1; i < len(phone); i++ {
		if phone[i] != first {
			return false
		}
	}
	return true
}
