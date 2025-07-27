package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/storage"
)

type FileService interface {
	UploadFile(file *multipart.FileHeader, description string, creatorID uint) (*models.File, error)
	GetFilesByCreator(creatorID uint) ([]*models.File, error)
	GetActiveByCreator(creatorID uint) ([]*models.File, error)
	GetFileByID(id uint) (*models.File, error)
	UpdateFile(id uint, description string) error
	DeleteFile(id uint) error
	GetFilesByType(creatorID uint, fileType string) ([]*models.File, error)
	ValidateFile(file *multipart.FileHeader) error
	GetFileType(ext string) string
}

type fileService struct {
	fileRepository repository.FileRepository
	s3Storage      storage.S3Storage
}

func NewFileService(fileRepository repository.FileRepository, s3Storage storage.S3Storage) FileService {
	return &fileService{
		fileRepository: fileRepository,
		s3Storage:      s3Storage,
	}
}

func (s *fileService) UploadFile(file *multipart.FileHeader, description string, creatorID uint) (*models.File, error) {
	// Validar arquivo
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	// Gerar nome único para o arquivo
	originalName := file.Filename
	fileExt := filepath.Ext(originalName)
	uniqueID := s.generateUniqueID()
	fileName := fmt.Sprintf("%s-%s%s",
		strings.TrimSuffix(originalName, fileExt),
		uniqueID,
		fileExt,
	)

	// Determinar tipo do arquivo
	fileType := s.getFileType(fileExt)

	// Upload para S3
	s3Key := fmt.Sprintf("files/%d/%s", creatorID, fileName)
	s3URL, err := s.s3Storage.UploadFile(file, s3Key)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload para S3: %w", err)
	}

	// Criar registro no banco
	fileModel := models.NewFile(
		fileName,
		originalName,
		description,
		fileType,
		s3Key,
		s3URL,
		file.Size,
		creatorID,
	)

	if err := s.fileRepository.Create(fileModel); err != nil {
		// Se falhar, tentar deletar do S3
		s.s3Storage.DeleteFile(s3Key)
		return nil, fmt.Errorf("erro ao salvar arquivo no banco: %w", err)
	}

	return fileModel, nil
}

func (s *fileService) GetFilesByCreator(creatorID uint) ([]*models.File, error) {
	return s.fileRepository.FindByCreator(creatorID)
}

func (s *fileService) GetActiveByCreator(creatorID uint) ([]*models.File, error) {
	return s.fileRepository.FindActiveByCreator(creatorID)
}

func (s *fileService) GetFileByID(id uint) (*models.File, error) {
	return s.fileRepository.FindByID(id)
}

func (s *fileService) UpdateFile(id uint, description string) error {
	file, err := s.fileRepository.FindByID(id)
	if err != nil {
		return err
	}

	file.Description = description
	return s.fileRepository.Update(file)
}

func (s *fileService) DeleteFile(id uint) error {
	file, err := s.fileRepository.FindByID(id)
	if err != nil {
		return err
	}

	// Deletar do S3
	if err := s.s3Storage.DeleteFile(file.S3Key); err != nil {
		return fmt.Errorf("erro ao deletar arquivo do S3: %w", err)
	}

	// Deletar do banco
	return s.fileRepository.Delete(id)
}

func (s *fileService) GetFilesByType(creatorID uint, fileType string) ([]*models.File, error) {
	return s.fileRepository.FindByType(creatorID, fileType)
}

func (s *fileService) ValidateFile(file *multipart.FileHeader) error {
	return s.validateFile(file)
}

func (s *fileService) GetFileType(ext string) string {
	return s.getFileType(ext)
}

func (s *fileService) validateFile(file *multipart.FileHeader) error {
	// Verificar tamanho (máximo 50MB)
	const maxSize = 50 * 1024 * 1024 // 50MB
	if file.Size > maxSize {
		return fmt.Errorf("arquivo muito grande. Tamanho máximo: 50MB")
	}

	// Verificar extensão
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png", ".gif"}

	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("tipo de arquivo não permitido. Tipos aceitos: %v", allowedExts)
	}

	return nil
}

func (s *fileService) getFileType(ext string) string {
	ext = strings.ToLower(ext)

	switch ext {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "document"
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	default:
		return "other"
	}
}

func (s *fileService) generateUniqueID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
