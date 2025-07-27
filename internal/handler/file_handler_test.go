package handler_test

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock FileService
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(file *multipart.FileHeader, description string, creatorID uint) (*models.File, error) {
	args := m.Called(file, description, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileService) GetFilesByCreator(creatorID uint) ([]*models.File, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func (m *MockFileService) GetActiveByCreator(creatorID uint) ([]*models.File, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func (m *MockFileService) GetFileByID(id uint) (*models.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileService) UpdateFile(id uint, description string) error {
	args := m.Called(id, description)
	return args.Error(0)
}

func (m *MockFileService) DeleteFile(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFileService) GetFilesByType(creatorID uint, fileType string) ([]*models.File, error) {
	args := m.Called(creatorID, fileType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.File), args.Error(1)
}

func (m *MockFileService) ValidateFile(file *multipart.FileHeader) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileService) GetFileType(ext string) string {
	args := m.Called(ext)
	return args.String(0)
}

func TestNewFileHandler(t *testing.T) {
	// Arrange
	mockFileService := &MockFileService{}
	mockSessionService := &mocks.MockSessionService{}

	// Act
	fileHandler := handler.NewFileHandler(mockFileService, mockSessionService)

	// Assert
	assert.NotNil(t, fileHandler)
}

func TestFileHandler_FileIndexView(t *testing.T) {
	// Arrange
	mockFileService := &MockFileService{}
	mockSessionService := &mocks.MockSessionService{}
	fileHandler := handler.NewFileHandler(mockFileService, mockSessionService)

	req, err := http.NewRequest("GET", "/file", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	fileHandler.FileIndexView(rr, req)

	// Assert
	assert.Equal(t, http.StatusSeeOther, rr.Code) // Deve redirecionar para login
}

func TestFileHandler_FileUploadView(t *testing.T) {
	// Arrange
	mockFileService := &MockFileService{}
	mockSessionService := &mocks.MockSessionService{}
	fileHandler := handler.NewFileHandler(mockFileService, mockSessionService)

	req, err := http.NewRequest("GET", "/file/upload", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	fileHandler.FileUploadView(rr, req)

	// Assert
	assert.Equal(t, http.StatusSeeOther, rr.Code) // Deve redirecionar para login
}

func TestFileHandler_FileDeleteSubmit(t *testing.T) {
	// Arrange
	mockFileService := &MockFileService{}
	mockSessionService := &mocks.MockSessionService{}
	fileHandler := handler.NewFileHandler(mockFileService, mockSessionService)

	req, err := http.NewRequest("POST", "/file/1/delete", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	fileHandler.FileDeleteSubmit(rr, req)

	// Assert
	assert.Equal(t, http.StatusSeeOther, rr.Code) // Deve redirecionar para login
}

func TestFileHandler_FileUpdateSubmit(t *testing.T) {
	// Arrange
	mockFileService := &MockFileService{}
	mockSessionService := &mocks.MockSessionService{}
	fileHandler := handler.NewFileHandler(mockFileService, mockSessionService)

	req, err := http.NewRequest("POST", "/file/1/update", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	fileHandler.FileUpdateSubmit(rr, req)

	// Assert
	assert.Equal(t, http.StatusSeeOther, rr.Code) // Deve redirecionar para login
}
