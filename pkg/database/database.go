package database

import (
	"log"
	"log/slog"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
	payment_models "github.com/anglesson/simple-web-server/internal/payment/data"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	if config.AppConfig.IsProduction() {
		connectWithPostgres()
	} else {
		connectWithSQLite()
	}
}

func connectWithPostgres() {
	connectGormAndMigrate(postgres.Open(config.AppConfig.DatabaseURL))
}

func connectWithSQLite() {
	connectGormAndMigrate(sqlite.Open("./mydb.db"))
}

func connectGormAndMigrate(dialector gorm.Dialector) {
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}
	migrate()
}

func migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Subscription{},
		&models.ClientCreator{},
		&models.Client{},
		&models.Contact{},
		&models.Creator{},
		&models.Ebook{},
		&models.Purchase{},
		&models.DownloadLog{},
		&payment_models.AccountModel{})

	if err != nil {
		slog.Error("Erro na geração das migrates")
		panic(err.Error())
	}
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Panic("failed to close database")
	}
	sqlDB.Close()
}
