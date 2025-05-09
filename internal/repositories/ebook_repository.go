package repositories

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"gorm.io/gorm"
)

type EbookQuery struct {
	Title       string
	Description string
	Pagination  *Pagination
}

type EbookRepository struct {
}

func NewEbookRepository() *EbookRepository {
	return &EbookRepository{}
}

func (r *EbookRepository) FindEbooksByUser(UserID uint, query EbookQuery) (*[]models.Ebook, error) {
	var ebooks []models.Ebook

	err := database.DB.
		Offset(query.Pagination.GetOffset()).Limit(query.Pagination.GetLimit()).
		Where("user_id = ?", UserID).
		InnerJoins("Creator").
		Scopes(ContainsTitleOrDescriptionWith(query.Title)).
		Find(&ebooks).
		Error
	if err != nil {
		log.Panicf("Erro na busca de ebooks: %s", err)
		return nil, errors.New("Erro na busca de dados")
	}

	query.Pagination.Total = int64(len(ebooks))

	return &ebooks, nil
}

func ContainsTitleOrDescriptionWith(term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title like '%" + term + "%'").Or("description like '%" + term + "%'")
	}
}
