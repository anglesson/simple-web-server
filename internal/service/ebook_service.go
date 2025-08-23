package service

import (
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/storage"
)

type EbookService interface {
	ListEbooksForUser(UserID string, query repository.EbookQuery) (*[]models.Ebook, error)
	FindByID(id uint) (*models.Ebook, error)
	FindBySlug(slug string) (*models.Ebook, error)
	Update(ebook *models.Ebook) error
	Create(ebook *models.Ebook) error
	Delete(id uint) error
}

type EbookServiceImpl struct {
	ebookRepository repository.EbookRepository
	s3Storage       storage.S3Storage
}

func NewEbookService(s3Storage storage.S3Storage) EbookService {
	return &EbookServiceImpl{
		ebookRepository: repository.NewGormEbookRepository(database.DB),
		s3Storage:       s3Storage,
	}
}

func (s *EbookServiceImpl) ListEbooksForUser(UserID string, query repository.EbookQuery) (*[]models.Ebook, error) {
	ebooks, err := s.ebookRepository.ListEbooksForUser(UserID, query)
	if err != nil {
		return nil, err
	}

	// Gerar URLs pré-assinadas para as imagens
	for i := range *ebooks {
		if (*ebooks)[i].Image != "" {
			(*ebooks)[i].Image = s.generatePresignedImageURL((*ebooks)[i].Image)
		}
	}

	return ebooks, nil
}

func (s *EbookServiceImpl) FindByID(id uint) (*models.Ebook, error) {
	ebook, err := s.ebookRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Gerar URL pré-assinada para a imagem
	if ebook.Image != "" {
		ebook.Image = s.generatePresignedImageURL(ebook.Image)
	}

	return ebook, nil
}

func (s *EbookServiceImpl) FindBySlug(slug string) (*models.Ebook, error) {
	ebook, err := s.ebookRepository.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Gerar URL pré-assinada para a imagem
	if ebook.Image != "" {
		ebook.Image = s.generatePresignedImageURL(ebook.Image)
	}

	return ebook, nil
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

// generatePresignedImageURL gera uma URL pré-assinada para a imagem de capa
func (s *EbookServiceImpl) generatePresignedImageURL(imageURL string) string {
	// Se a URL já é uma URL completa do S3, extrair a chave
	if imageURL == "" {
		return ""
	}

	// Para URLs do S3, extrair a chave do bucket
	// Exemplo: https://bucket.s3.region.amazonaws.com/ebook-covers/filename.jpg
	// Precisamos extrair: ebook-covers/filename.jpg

	// Se já é uma URL pré-assinada, retornar como está
	if len(imageURL) > 100 { // URLs pré-assinadas são muito longas
		return imageURL
	}

	// Se é uma URL pública do S3, gerar URL pré-assinada
	if s.isS3PublicURL(imageURL) {
		key := s.extractS3Key(imageURL)
		return s.s3Storage.GenerateDownloadLink(key)
	}

	// Se não é uma URL do S3, retornar como está (pode ser uma URL externa)
	return imageURL
}

// isS3PublicURL verifica se a URL é uma URL pública do S3
func (s *EbookServiceImpl) isS3PublicURL(url string) bool {
	return len(url) > 0 && (url[0:8] == "https://" || url[0:7] == "http://")
}

// extractS3Key extrai a chave S3 de uma URL pública
func (s *EbookServiceImpl) extractS3Key(url string) string {
	// Exemplo: https://bucket.s3.region.amazonaws.com/ebook-covers/filename.jpg
	// Retorna: ebook-covers/filename.jpg

	// Remover o protocolo
	if len(url) > 8 && url[0:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[0:7] == "http://" {
		url = url[7:]
	}

	// Para URLs do S3, o formato é: bucket.s3.region.amazonaws.com/ebook-covers/filename.jpg
	// Precisamos encontrar o terceiro '/' (após amazonaws.com)

	// Encontrar o primeiro '/' após o domínio
	firstSlash := -1
	for i, char := range url {
		if char == '/' {
			firstSlash = i
			break
		}
	}

	if firstSlash == -1 {
		return ""
	}

	// Para URLs do S3, a chave começa após o terceiro '/'
	// Exemplo: bucket.s3.region.amazonaws.com/ebook-covers/filename.jpg
	// Precisamos encontrar o '/' após "amazonaws.com"

	// Procurar por "amazonaws.com/"
	amazonawsIndex := strings.Index(url, "amazonaws.com/")
	if amazonawsIndex != -1 {
		// A chave começa após "amazonaws.com/"
		return url[amazonawsIndex+14:] // 14 = len("amazonaws.com/")
	}

	// Fallback: se não encontrar "amazonaws.com", usar a lógica anterior
	// Encontrar o segundo '/' (início da chave)
	secondSlash := -1
	for i := firstSlash + 1; i < len(url); i++ {
		if url[i] == '/' {
			secondSlash = i
			break
		}
	}

	if secondSlash == -1 {
		// Se não há segundo '/', a chave é tudo após o primeiro '/'
		return url[firstSlash+1:]
	}

	return url[secondSlash+1:]
}
