package database

import (
	"log"

	auth "github.com/anglesson/simple-web-server/internal/auth/models"
	ebook "github.com/anglesson/simple-web-server/internal/ebook/models"
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
	DB.AutoMigrate(&auth.User{})
	DB.AutoMigrate(&ebook.ClientCreator{})
	DB.AutoMigrate(&ebook.Client{})
	DB.AutoMigrate(&ebook.Contact{})
	DB.AutoMigrate(&ebook.Creator{})
	DB.AutoMigrate(&ebook.Ebook{})
	DB.AutoMigrate(&ebook.Purchase{})
}
