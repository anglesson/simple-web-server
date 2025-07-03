package domain

import (
	"fmt"
	"regexp"
	"strconv"
)

// CPF represents a Brazilian tax identification number (Cadastro de Pessoas Físicas)
type CPF struct {
	value string
}

// NewCPF creates a new CPF value object. Returns error if the CPF is invalid.
func NewCPF(value string) (*CPF, error) {
	cpf := &CPF{value: cleanCPF(value)}
	if !cpf.IsValid() {
		return nil, fmt.Errorf("invalid CPF: %s", value)
	}
	return cpf, nil
}

// Value returns the CPF value as a string
func (c *CPF) Value() string {
	return c.value
}

// String returns the formatted CPF (XXX.XXX.XXX-XX)
func (c *CPF) String() string {
	if len(c.value) != 11 {
		return c.value
	}
	return fmt.Sprintf("%s.%s.%s-%s",
		c.value[0:3],
		c.value[3:6],
		c.value[6:9],
		c.value[9:11],
	)
}

// IsValid checks if the CPF is valid according to Brazilian rules
func (c *CPF) IsValid() bool {
	// Check if CPF has 11 digits
	if len(c.value) != 11 {
		return false
	}

	// Check if all digits are the same (invalid CPF)
	if allDigitsSame(c.value) {
		return false
	}

	// Validate first digit
	digit1 := calculateDigit(c.value[:9], 10)
	if digit1 != int(c.value[9]-'0') {
		return false
	}

	// Validate second digit
	digit2 := calculateDigit(c.value[:10], 11)
	return digit2 == int(c.value[10]-'0')
}

// cleanCPF removes all non-digit characters from the CPF
func cleanCPF(cpf string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(cpf, "")
}

// allDigitsSame checks if all digits in the CPF are the same
func allDigitsSame(cpf string) bool {
	first := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

// calculateDigit calculates the verification digit for CPF
func calculateDigit(cpf string, factor int) int {
	var sum int
	for _, digit := range cpf {
		num, _ := strconv.Atoi(string(digit))
		sum += num * factor
		factor--
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

// Equal checks if two CPFs are equal
func (c *CPF) Equal(other *CPF) bool {
	if other == nil {
		return false
	}
	return c.value == other.value
}
