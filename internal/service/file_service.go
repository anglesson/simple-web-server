package service

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/storage"
)

type FileService interface {
	UploadFile(file *multipart.FileHeader, description string, creatorID uint) (*models.File, error)
	GetFilesByCreator(creatorID uint) ([]*models.File, error)
	GetFilesByCreatorPaginated(creatorID uint, query repository.FileQuery) ([]*models.File, int64, error)
	GetActiveByCreator(creatorID uint) ([]*models.File, error)
	GetFileByID(id uint) (*models.File, error)
	UpdateFile(id uint, name, description string) error
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
	files, err := s.fileRepository.FindByCreator(creatorID)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		file.S3URL = s.s3Storage.GenerateDownloadLink(file.S3Key)
	}
	return files, nil
}

func (s *fileService) GetFilesByCreatorPaginated(creatorID uint, query repository.FileQuery) ([]*models.File, int64, error) {
	query.CreatorID = creatorID
	files, total, err := s.fileRepository.FindByCreatorPaginated(query)
	if err != nil {
		return nil, 0, err
	}
	for _, file := range files {
		file.S3URL = s.s3Storage.GenerateDownloadLink(file.S3Key)
	}
	return files, total, nil
}

func (s *fileService) GetActiveByCreator(creatorID uint) ([]*models.File, error) {
	return s.fileRepository.FindActiveByCreator(creatorID)
}

func (s *fileService) GetFileByID(id uint) (*models.File, error) {
	return s.fileRepository.FindByID(id)
}

func (s *fileService) UpdateFile(id uint, name, description string) error {
	file, err := s.fileRepository.FindByID(id)
	if err != nil {
		return err
	}

	file.Name = name
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
	// Verificar tamanho (máximo 10MB - reduzido por segurança)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxSize {
		return fmt.Errorf("arquivo muito grande. Tamanho máximo: 10MB")
	}

	// Verificar se o arquivo não está vazio
	if file.Size == 0 {
		return fmt.Errorf("arquivo vazio não é permitido")
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

	// Verificar MIME type
	if err := s.validateMimeType(file); err != nil {
		return err
	}

	// Verificar conteúdo do arquivo para detectar arquivos maliciosos
	if err := s.validateFileContent(file); err != nil {
		return err
	}

	return nil
}

// validateMimeType validates the actual MIME type of the file
func (s *fileService) validateMimeType(file *multipart.FileHeader) error {
	// Open file to check MIME type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer src.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Detect MIME type
	mimeType := http.DetectContentType(buffer)

	// Allowed MIME types
	allowedMimeTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("tipo MIME não permitido: %s", mimeType)
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

// validateFileContent checks for malicious content in files
func (s *fileService) validateFileContent(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo para validação: %w", err)
	}
	defer src.Close()

	// Read first 1024 bytes to check for malicious signatures
	buffer := make([]byte, 1024)
	bytesRead, err := src.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("erro ao ler arquivo para validação: %w", err)
	}

	if bytesRead == 0 {
		return fmt.Errorf("arquivo vazio")
	}

	// Check for common malicious file signatures
	maliciousSignatures := [][]byte{
		{0x4D, 0x5A}, // MZ header (executables)
		{0x7F, 0x45, 0x4C, 0x46}, // ELF header
		{0xFE, 0xED, 0xFA, 0xCE}, // Mach-O header
		{0x50, 0x4B, 0x03, 0x04}, // ZIP with potential executable content
		{0x1F, 0x8B, 0x08}, // GZIP
		{0x25, 0x50, 0x44, 0x46}, // PDF (but check for embedded scripts)
	}

	for _, signature := range maliciousSignatures {
		if bytesRead >= len(signature) && bytes.Equal(buffer[:len(signature)], signature) {
			// For PDFs, we need additional checks
			if bytes.Equal(signature, []byte{0x25, 0x50, 0x44, 0x46}) {
				// Check if PDF contains JavaScript
				if bytes.Contains(buffer, []byte("/JS")) || bytes.Contains(buffer, []byte("/JavaScript")) {
					return fmt.Errorf("PDF contém JavaScript que pode ser malicioso")
				}
			} else {
				return fmt.Errorf("tipo de arquivo potencialmente perigoso detectado")
			}
		}
	}

	// Check for script content in text files
	textContent := strings.ToLower(string(buffer))
	scriptKeywords := []string{
		"<script", "javascript:", "vbscript:", "onload=", "onerror=",
		"eval(", "document.cookie", "window.location", "alert(",
		"<?php", "<?=", "<%", "%>", "<?", "?>",
	}

	for _, keyword := range scriptKeywords {
		if strings.Contains(textContent, keyword) {
			return fmt.Errorf("arquivo contém conteúdo de script potencialmente perigoso")
		}
	}

	return nil
}

func (s *fileService) generateUniqueID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
