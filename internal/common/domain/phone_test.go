package common_domain_test

import (
	"testing"

	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
)

func TestNewPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+1234567890", true},
		{"1234567890", true},
		{"invalid-phone", false},
		{"", false},
	}

	for _, test := range tests {
		_, err := common_domain.NewPhone(test.input)
		if (err == nil) != test.expected {
			t.Errorf("NewPhone(%q) expected %v, got %v", test.input, test.expected, err == nil)
		}
	}
}
