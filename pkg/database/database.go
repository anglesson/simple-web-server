package database

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	if config.AppConfig.IsProdcution() {
		DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Panic("failed to connect database")
		}
	} else {
		DB, err = gorm.Open(sqlite.Open("./mydb.db"), &gorm.Config{})
		if err != nil {
			log.Panic("failed to connect database")
		}
	}

	migrate()
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
