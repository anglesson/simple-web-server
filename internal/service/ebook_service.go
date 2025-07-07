package service

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
)

type EbookService struct{}

func NewEbookService() *EbookService {
	return &EbookService{}
}

func (s *EbookService) ListEbooksForUser(UserID uint, query repository.EbookQuery) (*[]models.Ebook, error) {
	ebookRepository := repository.NewEbookRepository()
	ebooks, err := ebookRepository.FindEbooksByUser(UserID, query)
	if err != nil {
		return nil, errors.New("ebooks não encontrados")
	}

	if len(*ebooks) == 0 {
		return nil, nil
	}

	return ebooks, nil
}
