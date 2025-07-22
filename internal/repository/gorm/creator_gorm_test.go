package gorm_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/pkg/database"
)

func TestCreatorGorm_Save(t *testing.T) {
	database.Connect()

	tx := database.DB.Begin()
	defer tx.Rollback()

	// Arrange
	creator, err := domain.NewCreator("valid name", "valid@mail.com", "058.997.950-77", "(81) 98765-4321", "1990-01-01")
	if err != nil {
		t.Helper()
		t.Fatal(err)
	}
	sut := gorm.NewCreatorRepository(database.DB)

	// Act
	err = sut.Save(creator)
	if err != nil {
		t.Fatalf("unexpected error saving creator: %v", err)
	}

	// Assert
	if creator.ID == 0 {
		t.Errorf("expected creator ID to be greater than 0, got %d", creator.ID)
	}

	// Clean up (optional, if your repository supports Delete)
	// _ = sut.Delete(creator.ID)
}
