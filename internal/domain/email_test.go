package domain_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewEmail_ValidEmail(t *testing.T) {
	email, err := domain.NewEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Value())
}

func TestNewEmail_InvalidEmail(t *testing.T) {
	_, err := domain.NewEmail("invalid-email")
	assert.Error(t, err)
	assert.Equal(t, "endereço de email inválido", err.Error())
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := domain.NewEmail("test@example.com")
	email2, _ := domain.NewEmail("test@example.com")
	email3, _ := domain.NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

func TestEmail_String(t *testing.T) {
	email, _ := domain.NewEmail("test@example.com")
	assert.Equal(t, "test@example.com", email.String())
}
