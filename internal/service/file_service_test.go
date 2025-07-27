package service_test

import (
	"mime/multipart"
	"testing"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock S3Storage
type MockS3Storage struct {
	mock.Mock
}

func (m *MockS3Storage) UploadFile(file *multipart.FileHeader, key string) (string, error) {
	args := m.Called(file, key)
	return args.String(0), args.Error(1)
}

func (m *MockS3Storage) DeleteFile(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockS3Storage) GenerateDownloadLink(key string) string {
	args := m.Called(key)
	return args.String(0)
}

// Mock FileRepository
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) FindByID(id uint) (*models.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileRepository) FindByCreator(creatorID uint) ([]*models.File, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func (m *MockFileRepository) FindActiveByCreator(creatorID uint) ([]*models.File, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func (m *MockFileRepository) Update(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFileRepository) FindByType(creatorID uint, fileType string) ([]*models.File, error) {
	args := m.Called(creatorID, fileType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func TestNewFileService(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}

	// Act
	fileService := service.NewFileService(mockRepo, mockStorage)

	// Assert
	assert.NotNil(t, fileService)
}

func TestFileService_GetFilesByCreator(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	creatorID := uint(1)
	expectedFiles := []*models.File{
		{Name: "file1.pdf", CreatorID: creatorID},
		{Name: "file2.pdf", CreatorID: creatorID},
	}

	mockRepo.On("FindByCreator", creatorID).Return(expectedFiles, nil)
	mockStorage.On("GenerateDownloadLink", mock.Anything).Return("https://dummy-presigned-url", nil)

	// Act
	files, err := fileService.GetFilesByCreator(creatorID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, expectedFiles, files)
	mockRepo.AssertExpectations(t)
}

func TestFileService_GetFileByID(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	fileID := uint(1)
	expectedFile := &models.File{Name: "test.pdf"}

	mockRepo.On("FindByID", fileID).Return(expectedFile, nil)

	// Act
	file, err := fileService.GetFileByID(fileID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFile, file)
	mockRepo.AssertExpectations(t)
}

func TestFileService_UpdateFile(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	fileID := uint(1)
	description := "Updated description"
	existingFile := &models.File{Description: "Old description"}

	mockRepo.On("FindByID", fileID).Return(existingFile, nil)
	mockRepo.On("Update", existingFile).Return(nil)

	// Act
	err := fileService.UpdateFile(fileID, description)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, description, existingFile.Description)
	mockRepo.AssertExpectations(t)
}

func TestFileService_DeleteFile(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	fileID := uint(1)
	existingFile := &models.File{S3Key: "files/1/test.pdf"}

	mockRepo.On("FindByID", fileID).Return(existingFile, nil)
	mockStorage.On("DeleteFile", existingFile.S3Key).Return(nil)
	mockRepo.On("Delete", fileID).Return(nil)

	// Act
	err := fileService.DeleteFile(fileID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestFileService_GetFilesByType(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	creatorID := uint(1)
	fileType := "pdf"
	expectedFiles := []*models.File{
		{Name: "file1.pdf", FileType: "pdf", CreatorID: creatorID},
	}

	mockRepo.On("FindByType", creatorID, fileType).Return(expectedFiles, nil)

	// Act
	files, err := fileService.GetFilesByType(creatorID, fileType)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, expectedFiles, files)
	mockRepo.AssertExpectations(t)
}

func TestFileService_GetActiveByCreator(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	creatorID := uint(1)
	expectedFiles := []*models.File{
		{Name: "active1.pdf", Status: true, CreatorID: creatorID},
		{Name: "active2.pdf", Status: true, CreatorID: creatorID},
	}

	mockRepo.On("FindActiveByCreator", creatorID).Return(expectedFiles, nil)

	// Act
	files, err := fileService.GetActiveByCreator(creatorID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, expectedFiles, files)
	mockRepo.AssertExpectations(t)
}

// Testes para validação de arquivos
func TestFileService_validateFile(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	tests := []struct {
		name        string
		filename    string
		size        int64
		expectError bool
	}{
		{
			name:        "Valid PDF file",
			filename:    "test.pdf",
			size:        1024 * 1024, // 1MB
			expectError: false,
		},
		{
			name:        "Valid image file",
			filename:    "image.jpg",
			size:        512 * 1024, // 512KB
			expectError: false,
		},
		{
			name:        "File too large",
			filename:    "large.pdf",
			size:        100 * 1024 * 1024, // 100MB
			expectError: true,
		},
		{
			name:        "Invalid file type",
			filename:    "script.exe",
			size:        1024 * 1024,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar um mock de FileHeader
			fileHeader := &multipart.FileHeader{
				Filename: tt.filename,
				Size:     tt.size,
			}

			// Act
			err := fileService.ValidateFile(fileHeader)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileService_getFileType(t *testing.T) {
	// Arrange
	mockRepo := &MockFileRepository{}
	mockStorage := &MockS3Storage{}
	fileService := service.NewFileService(mockRepo, mockStorage)

	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{name: ".pdf", ext: ".pdf", expected: "pdf"},
		{name: ".doc", ext: ".doc", expected: "document"},
		{name: ".docx", ext: ".docx", expected: "document"},
		{name: ".jpg", ext: ".jpg", expected: "image"},
		{name: ".jpeg", ext: ".jpeg", expected: "image"},
		{name: ".png", ext: ".png", expected: "image"},
		{name: ".gif", ext: ".gif", expected: "image"},
		{name: ".txt", ext: ".txt", expected: "other"},
		{name: ".zip", ext: ".zip", expected: "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := fileService.GetFileType(tt.name)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
