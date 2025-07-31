package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mime/multipart"

	handler "github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	service_mocks "github.com/anglesson/simple-web-server/internal/service/mocks"
	template_mocks "github.com/anglesson/simple-web-server/pkg/template/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock EbookService for testing
type MockEbookService struct {
	mock.Mock
}

func (m *MockEbookService) ListEbooksForUser(UserID uint, query repository.EbookQuery) (*[]models.Ebook, error) {
	args := m.Called(UserID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.Ebook), args.Error(1)
}

func (m *MockEbookService) FindByID(id uint) (*models.Ebook, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ebook), args.Error(1)
}

func (m *MockEbookService) FindBySlug(slug string) (*models.Ebook, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ebook), args.Error(1)
}

func (m *MockEbookService) Update(ebook *models.Ebook) error {
	args := m.Called(ebook)
	return args.Error(0)
}

func (m *MockEbookService) Create(ebook *models.Ebook) error {
	args := m.Called(ebook)
	return args.Error(0)
}

func (m *MockEbookService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock S3Storage for testing
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

func (m *MockS3Storage) GenerateDownloadLinkWithExpiration(key string, expirationSeconds int) string {
	args := m.Called(key, expirationSeconds)
	return args.String(0)
}

// Mock FlashMessage for testing
type MockFlashMessage struct {
	mock.Mock
}

func (m *MockFlashMessage) Success(message string) {
	m.Called(message)
}

func (m *MockFlashMessage) Error(message string) {
	m.Called(message)
}

func (m *MockFlashMessage) Warning(message string) {
	m.Called(message)
}

func (m *MockFlashMessage) Info(message string) {
	m.Called(message)
}

// Mock CreatorService for testing
type MockCreatorService struct {
	mock.Mock
}

func (m *MockCreatorService) CreateCreator(input service.InputCreateCreator) (*models.Creator, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorService) FindCreatorByEmail(email string) (*models.Creator, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorService) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Creator), args.Error(1)
}

func (m *MockCreatorService) FindByID(id uint) (*models.Creator, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Creator), args.Error(1)
}

type EbookHandlerTestSuite struct {
	suite.Suite
	sut                  *handler.EbookHandler
	mockEbookService     *MockEbookService
	mockCreatorService   *MockCreatorService
	mockFileService      *service_mocks.MockFileService
	mockS3Storage        *MockS3Storage
	mockFlashMessage     *MockFlashMessage
	mockTemplateRenderer *template_mocks.MockTemplateRenderer
	flashFactory         web.FlashMessageFactory
}

func (suite *EbookHandlerTestSuite) SetupTest() {
	suite.mockEbookService = new(MockEbookService)
	suite.mockCreatorService = new(MockCreatorService)
	suite.mockFileService = new(service_mocks.MockFileService)
	suite.mockS3Storage = new(MockS3Storage)
	suite.mockFlashMessage = new(MockFlashMessage)
	suite.mockTemplateRenderer = new(template_mocks.MockTemplateRenderer)

	suite.flashFactory = func(w http.ResponseWriter, r *http.Request) web.FlashMessagePort {
		return suite.mockFlashMessage
	}

	suite.sut = handler.NewEbookHandler(
		suite.mockEbookService,
		suite.mockCreatorService,
		suite.mockFileService,
		suite.mockS3Storage,
		suite.flashFactory,
		suite.mockTemplateRenderer,
	)
}

// Helper function to create a request with mocked user
func (suite *EbookHandlerTestSuite) createRequestWithUser(email string) *http.Request {
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), middleware.UserEmailKey, email)
	return r.WithContext(ctx)
}

func (suite *EbookHandlerTestSuite) TestIndexView_UserNotFound() {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	// No user email in context

	// Act
	suite.sut.IndexView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *EbookHandlerTestSuite) TestIndexView_ServiceError() {
	// Arrange
	w := httptest.NewRecorder()
	r := suite.createRequestWithUser("test@example.com")

	suite.mockEbookService.On("ListEbooksForUser", uint(1), mock.AnythingOfType("repository.EbookQuery")).
		Return(nil, errors.New("service error"))

	// Act
	suite.sut.IndexView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockEbookService.AssertExpectations(suite.T())
}

func (suite *EbookHandlerTestSuite) TestCreateView_UserNotFound() {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	// No user email in context

	// Act
	suite.sut.CreateView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *EbookHandlerTestSuite) TestCreateView_CreatorNotFound() {
	// Arrange
	w := httptest.NewRecorder()
	r := suite.createRequestWithUser("test@example.com")

	suite.mockCreatorService.On("FindCreatorByUserID", uint(1)).
		Return(nil, errors.New("creator not found"))

	// Act
	suite.sut.CreateView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockCreatorService.AssertExpectations(suite.T())
}

func (suite *EbookHandlerTestSuite) TestShowView_UserNotFound() {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	// No user email in context

	// Act
	suite.sut.ShowView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *EbookHandlerTestSuite) TestShowView_EbookNotFound() {
	// Arrange
	w := httptest.NewRecorder()
	r := suite.createRequestWithUser("test@example.com")

	// Mock chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	suite.mockEbookService.On("FindByID", uint(1)).
		Return(nil, errors.New("ebook not found"))

	// Act
	suite.sut.ShowView(w, r)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockEbookService.AssertExpectations(suite.T())
}

func TestEbookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(EbookHandlerTestSuite))
}
