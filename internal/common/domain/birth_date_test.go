package common_domain_test

import (
	"testing"
	"time"

	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewBirthDate_ValidDate(t *testing.T) {
	birthDate, err := common_domain.NewBirthDate("2000-01-01")
	assert.NoError(t, err)
	assert.Equal(t, "2000-01-01", birthDate.String())
}

func TestNewBirthDate_InvalidFormat(t *testing.T) {
	_, err := common_domain.NewBirthDate("01-01-2000")
	assert.Error(t, err)
	assert.Equal(t, "data de nascimento inválida: formato deve ser YYYY-MM-DD", err.Error())
}

func TestNewBirthDate_FutureDate(t *testing.T) {
	futureDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	_, err := common_domain.NewBirthDate(futureDate)
	assert.Error(t, err)
	assert.Equal(t, "data de nascimento não pode ser no futuro", err.Error())
}

func TestBirthDate_Value(t *testing.T) {
	birthDate, _ := common_domain.NewBirthDate("2000-01-01")
	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), birthDate.Value())
}

func TestBirthDate_String(t *testing.T) {
	birthDate, _ := common_domain.NewBirthDate("2000-01-01")
	assert.Equal(t, "2000-01-01", birthDate.String())
}
