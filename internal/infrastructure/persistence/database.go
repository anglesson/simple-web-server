package persistence

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/shared/database"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = database.GetDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func GetDB() *gorm.DB {
	return db
}
