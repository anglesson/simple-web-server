package database

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./mydb.db"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	migrate()
}

func migrate() {
	DB.AutoMigrate(&models.User{})
}
