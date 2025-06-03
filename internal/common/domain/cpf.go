package common_domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type CPF struct {
	value string
}

func NewCPF(number string) (CPF, error) {
	// Remover pontuações e espaços
	cleanedNumber := regexp.MustCompile(`[\.\-/\s]`).ReplaceAllString(number, "")

	if len(cleanedNumber) != 11 {
		return CPF{}, errors.New("CPF deve ter 11 dígitos")
	}

	if !isValidCPFChecksum(cleanedNumber) { // Implementar a validação de dígitos verificadores
		return CPF{}, errors.New("CPF inválido (dígitos verificadores)")
	}

	return CPF{value: cleanedNumber}, nil
}

func (c CPF) Value() string {
	return c.value
}

func (c CPF) Equals(other CPF) bool {
	return c.value == other.value
}

func (c CPF) String() string {
	return fmt.Sprintf("%s.%s.%s-%s", c.value[:3], c.value[3:6], c.value[6:9], c.value[9:])
}

// isValidCPFChecksum é uma função auxiliar para validar os dígitos verificadores do CPF.
func isValidCPFChecksum(cpf string) bool {
	// Validate format: CPF must be 11 digits and match a numeric pattern.
	if !regexp.MustCompile(`^\d{11}$`).MatchString(cpf) {
		return false
	}

	// Check for invalid sequences (e.g., all digits the same).
	invalidSequences := []string{
		"00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999",
	}
	for _, seq := range invalidSequences {
		if cpf == seq {
			return false
		}
	}

	// Validate checksum digits.
	digits := make([]int, 11)
	for i := 0; i < 11; i++ {
		digit, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Calculate first checksum digit.
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	firstCheck := (sum * 10) % 11
	if firstCheck == 10 {
		firstCheck = 0
	}
	if digits[9] != firstCheck {
		return false
	}

	// Calculate second checksum digit.
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	secondCheck := (sum * 10) % 11
	if secondCheck == 10 {
		secondCheck = 0
	}
	return digits[10] == secondCheck
}
