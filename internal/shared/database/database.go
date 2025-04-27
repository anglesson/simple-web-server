package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./mydb.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetDB() *sql.DB {
	if db == nil {
		db = Connect()
	}
	return db
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Migrate() {
	// Migrate the database schema
	db = GetDB()
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT UNIQUE,
		password TEXT
	);
	`)
	if err != nil {
		log.Fatal("Erro ao gerar migrations ", err)
	}
}
