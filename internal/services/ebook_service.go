package services

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
)

type EbookService struct{}

func NewEbookService() *EbookService {
	return &EbookService{}
}

func (s *EbookService) ListEbooksForUser(UserID uint, pagination *repositories.Pagination) (*[]models.Ebook, error) {
	ebookRepository := repositories.NewEbookRepository()
	ebooks, err := ebookRepository.FindEbooksByUser(UserID, pagination)
	if err != nil {
		return nil, errors.New("ebooks não encontrados")
	}

	if len(*ebooks) == 0 {
		return nil, nil
	}

	return ebooks, nil
}
