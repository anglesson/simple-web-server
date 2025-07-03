package repositories

import (
	"errors"
	common_application "github.com/anglesson/simple-web-server/domain"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"gorm.io/gorm"
)

type EbookQuery struct {
	Title       string
	Description string
	Pagination  *common_application.Pagination
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
		Offset(query.Pagination.GetOffset()).
		Limit(query.Pagination.GetLimit()).
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
		return db.Where("title LIKE ? OR description LIKE ?", "%"+term+"%", "%"+term+"%")
	}
}
