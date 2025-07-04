package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	// In a real application, use a secure hashing algorithm
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}

	// Convert the hashed password to a string
	hashedPassword := string(bytes)

	// Return the hashed password
	return hashedPassword
}

func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// Error response structure
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ServerError handles 500 internal server errors
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("SERVER ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	http.ServeFile(w, r, "internal/templates/pages/500.html")
}

// ClientError handles 4xx client errors
func ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("CLIENT ERROR: %s", message)
	w.WriteHeader(status)

	switch status {
	case http.StatusNotFound:
		http.ServeFile(w, r, "internal/templates/pages/404.html")
	case http.StatusUnauthorized:
		http.ServeFile(w, r, "internal/templates/pages/unauthorized.html")
	default:
		http.ServeFile(w, r, "internal/templates/pages/error.html")
	}
}

// NotFound returns a 404 not found error
func NotFound(w http.ResponseWriter, r *http.Request) {
	ClientError(w, r, http.StatusNotFound, "The requested resource could not be found")
}
