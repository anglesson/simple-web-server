package repositories

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

func Save(user *models.User) {
	db := database.GetDB()

	insertStmt := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := db.Exec(insertStmt, user.Username, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Printf("Created new user with %s", user.Email)

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("LAST ID: %d", id)
	findByID(id, user)
}

func findByID(id int64, user *models.User) {
	db := database.GetDB()

	selectStmt := `SELECT id, username, email FROM users WHERE id = ?`
	row := db.QueryRow(selectStmt, id)
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		log.Fatalf("Error to scan user: %s", err)
	}
}

func FindByEmail(emailUser string) *models.User {
	db := database.GetDB()

	user := &models.User{}

	selectStmt := `SELECT id, username, email FROM users WHERE email = ?`
	row := db.QueryRow(selectStmt, emailUser)

	row.Scan(&user.ID, &user.Username, &user.Email)

	return user
}
