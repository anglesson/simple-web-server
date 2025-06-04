package common_domain_test

import (
	"testing"

	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewEmail_ValidEmail(t *testing.T) {
	email, err := common_domain.NewEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Value())
}

func TestNewEmail_InvalidEmail(t *testing.T) {
	_, err := common_domain.NewEmail("invalid-email")
	assert.Error(t, err)
	assert.Equal(t, "endereço de email inválido", err.Error())
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := common_domain.NewEmail("test@example.com")
	email2, _ := common_domain.NewEmail("test@example.com")
	email3, _ := common_domain.NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

func TestEmail_String(t *testing.T) {
	email, _ := common_domain.NewEmail("test@example.com")
	assert.Equal(t, "test@example.com", email.String())
}
