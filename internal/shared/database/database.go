package database

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
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

func GetDB() (*gorm.DB, error) {
	if DB == nil {
		Connect()
	}
	return DB, nil
}

func migrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.ClientCreator{})
	DB.AutoMigrate(&models.Client{})
	DB.AutoMigrate(&models.Contact{})
	DB.AutoMigrate(&models.Creator{})
	DB.AutoMigrate(&models.Ebook{})
	DB.AutoMigrate(&models.Purchase{})
	DB.AutoMigrate(&models.DownloadLog{})
}
