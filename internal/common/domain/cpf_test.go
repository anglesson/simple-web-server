package common_domain_test

import (
	"testing"

	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewCPF_ValidCPF(t *testing.T) {
	cpf, err := common_domain.NewCPF("12345678909")
	assert.NoError(t, err)
	assert.Equal(t, "12345678909", cpf.Value())
}

func TestNewCPF_InvalidCPF(t *testing.T) {
	_, err := common_domain.NewCPF("12345678900")
	assert.Error(t, err)
	assert.Equal(t, "CPF inválido (dígitos verificadores)", err.Error())
}

func TestNewCPF_InvalidLength(t *testing.T) {
	_, err := common_domain.NewCPF("12345")
	assert.Error(t, err)
	assert.Equal(t, "CPF deve ter 11 dígitos", err.Error())
}

func TestCPF_Equals(t *testing.T) {
	cpf1, _ := common_domain.NewCPF("12345678909")
	cpf2, _ := common_domain.NewCPF("12345678909")
	cpf3, _ := common_domain.NewCPF("98765432100")

	assert.True(t, cpf1.Equals(cpf2))
	assert.False(t, cpf1.Equals(cpf3))
}

func TestCPF_String(t *testing.T) {
	cpf, _ := common_domain.NewCPF("12345678909")
	assert.Equal(t, "123.456.789-09", cpf.String())
}
