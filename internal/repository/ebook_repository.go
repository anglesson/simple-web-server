package repository

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/database"
	"gorm.io/gorm"
)

type EbookQuery struct {
	Title       string
	Description string
	Pagination  *models.Pagination
}

type EbookRepository struct {
}

func NewEbookRepository() *EbookRepository {
	return &EbookRepository{}
}

func (r *EbookRepository) FindEbooksByUser(UserID uint, query EbookQuery) (*[]models.Ebook, error) {
	var ebooks []models.Ebook
	log.Printf("UsuarioID: %v", UserID)
	err := database.DB.
		Model(&models.Ebook{}).
		Joins("INNER JOIN creators ON creators.id = ebooks.creator_id").
		Where("creators.user_id = ?", UserID).
		Scopes(ContainsTitleOrDescriptionWith(query.Title)).
		Offset(getOffset(query.Pagination)).
		Limit(getLimit(query.Pagination)).
		Find(&ebooks).
		Error
	if err != nil {
		log.Panicf("Erro na busca de ebooks: %s", err)
		return nil, errors.New("Erro na busca de dados")
	}

	return &ebooks, nil
}

func ContainsTitleOrDescriptionWith(term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title LIKE ? OR description LIKE ?", "%"+term+"%", "%"+term+"%")
	}
}

// Helper functions for pagination
func getOffset(pagination *models.Pagination) int {
	if pagination == nil {
		return 0
	}
	return (pagination.Page - 1) * pagination.Limit
}

func getLimit(pagination *models.Pagination) int {
	if pagination == nil {
		return 10 // default limit
	}
	return pagination.Limit
}
