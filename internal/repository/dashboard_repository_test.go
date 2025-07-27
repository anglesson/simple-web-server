package repository_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory sqlite: %v", err)
	}
	return db
}

func setupDashboardTestDB(t *testing.T) *gorm.DB {
	db := testDB(t)
	db.AutoMigrate(&models.Client{}, &models.Purchase{}, &models.Ebook{}, &models.Creator{})
	return db
}

func TestGetTopClients_Success(t *testing.T) {
	db := setupDashboardTestDB(t)
	dr := repository.NewDashboardRepository(1)
	database.DB = db

	creator := models.Creator{Model: gorm.Model{ID: 1}, Name: "Creator1", UserID: 1}
	db.Create(&creator)
	client := models.Client{Name: "Client1", Email: "client1@email.com"}
	db.Create(&client)
	ebook := models.Ebook{Title: "Ebook1", CreatorID: creator.ID}
	db.Create(&ebook)
	purchase := models.Purchase{ClientID: client.ID, EbookID: ebook.ID}
	db.Create(&purchase)

	topClients, err := dr.GetTopClients()
	assert.NoError(t, err)
	assert.Len(t, topClients, 1)
	assert.Equal(t, "Client1", topClients[0].Name)
	assert.Equal(t, "client1@email.com", topClients[0].Email)
	assert.Equal(t, int64(1), topClients[0].TotalPurchases)
}

func TestGetTopClients_NoClients(t *testing.T) {
	db := setupDashboardTestDB(t)
	dr := repository.NewDashboardRepository(1)
	database.DB = db

	topClients, err := dr.GetTopClients()
	assert.NoError(t, err)
	assert.Len(t, topClients, 0)
}
