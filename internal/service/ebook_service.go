package service

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
)

type EbookService struct {
	ebookRepository repository.EbookRepository
}

func NewEbookService() *EbookService {
	return &EbookService{
		ebookRepository: repository.NewGormEbookRepository(database.DB),
	}
}

func (s *EbookService) ListEbooksForUser(UserID uint, query repository.EbookQuery) (*[]models.Ebook, error) {
	return s.ebookRepository.ListEbooksForUser(UserID, query)
}
