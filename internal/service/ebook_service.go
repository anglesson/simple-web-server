package service

import (
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
)

type EbookService interface {
	ListEbooksForUser(UserID uint, query repository.EbookQuery) (*[]models.Ebook, error)
	FindByID(id uint) (*models.Ebook, error)
	FindBySlug(slug string) (*models.Ebook, error)
	Update(ebook *models.Ebook) error
	Create(ebook *models.Ebook) error
	Delete(id uint) error
}

type EbookServiceImpl struct {
	ebookRepository repository.EbookRepository
}

func NewEbookService() EbookService {
	return &EbookServiceImpl{
		ebookRepository: repository.NewGormEbookRepository(database.DB),
	}
}

func (s *EbookServiceImpl) ListEbooksForUser(UserID uint, query repository.EbookQuery) (*[]models.Ebook, error) {
	return s.ebookRepository.ListEbooksForUser(UserID, query)
}

func (s *EbookServiceImpl) FindByID(id uint) (*models.Ebook, error) {
	return s.ebookRepository.FindByID(id)
}

func (s *EbookServiceImpl) FindBySlug(slug string) (*models.Ebook, error) {
	return s.ebookRepository.FindBySlug(slug)
}

func (s *EbookServiceImpl) Update(ebook *models.Ebook) error {
	return s.ebookRepository.Update(ebook)
}

func (s *EbookServiceImpl) Create(ebook *models.Ebook) error {
	return s.ebookRepository.Create(ebook)
}

func (s *EbookServiceImpl) Delete(id uint) error {
	return s.ebookRepository.Delete(id)
}
