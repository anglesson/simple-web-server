package gorm_test

import (
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/pkg/database"
)

func TestCreatorGorm_Create(t *testing.T) {
	database.Connect()

	tx := database.DB.Begin()
	defer tx.Rollback()

	// Arrange
	birthDate, _ := time.Parse("2006-01-02", "1990-01-01")
	creator := models.NewCreator(
		"valid name",
		"valid@mail.com",
		"81987654321",
		"05899795077",
		birthDate,
		1, // userID
	)
	sut := gorm.NewCreatorRepository(database.DB)

	// Act
	err := sut.Create(creator)
	if err != nil {
		t.Fatalf("unexpected error creating creator: %v", err)
	}

	// Assert
	if creator.ID == 0 {
		t.Errorf("expected creator ID to be greater than 0, got %d", creator.ID)
	}
}
