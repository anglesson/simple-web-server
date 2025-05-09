package repositories

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

type EbookRepository struct {
}

func NewEbookRepository() *EbookRepository {
	return &EbookRepository{}
}

func (r *EbookRepository) FindEbooksByUser(UserID uint, pagination *Pagination) (*[]models.Ebook, error) {
	var ebooks []models.Ebook

	err := database.DB.
		Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).
		Where("user_id = ?", UserID).
		InnerJoins("Creator").
		Find(&ebooks).
		Error
	if err != nil {
		log.Panicf("Erro na busca de ebooks: %s", err)
		return nil, errors.New("Erro na busca de dados")
	}

	pagination.Total = int64(len(ebooks))

	return &ebooks, nil
}
