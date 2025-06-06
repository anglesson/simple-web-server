package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/domain"
	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server/handlers"
	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockClientUseCase struct {
	mock.Mock
}

func (m *MockClientUseCase) CreateClient(input application.CreateClientInput) (*application.CreateClientOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.CreateClientOutput), args.Error(1)
}

func (m *MockClientUseCase) UpdateClient(input application.UpdateClientInput) (*application.UpdateClientOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.UpdateClientOutput), args.Error(1)
}

func (m *MockClientUseCase) ImportClients(input application.ImportClientsInput) (*application.ImportClientsOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.ImportClientsOutput), args.Error(1)
}

func (m *MockClientUseCase) ListClients(input application.ListClientsInput) (*application.ListClientsOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.ListClientsOutput), args.Error(1)
}

type ClientHandlerTestSuite struct {
	suite.Suite
	mockUseCase *MockClientUseCase
	handler     *handlers.ClientHandler
	tmpDir      string
}

func (s *ClientHandlerTestSuite) SetupSuite() {
	// Create a temporary directory for test templates
	var err error
	s.tmpDir, err = os.MkdirTemp("", "templates")
	s.NoError(err)

	// Create a test template file
	templateContent := `<!DOCTYPE html>
<html>
<body>
    {{if .Clients}}
        {{range .Clients}}
            <div>{{.Name}}</div>
        {{end}}
    {{else}}
        <div>No clients found</div>
    {{end}}
    {{if .Pagination}}
        <div>Pagination: {{.Pagination}}</div>
    {{end}}
</body>
</html>`

	// Create the client directory
	clientDir := filepath.Join(s.tmpDir, "client")
	err = os.MkdirAll(clientDir, 0755)
	s.NoError(err)

	// Write the template file
	err = os.WriteFile(filepath.Join(clientDir, "index.html"), []byte(templateContent), 0644)
	s.NoError(err)
}

func (s *ClientHandlerTestSuite) TearDownSuite() {
	os.RemoveAll(s.tmpDir)
}

func (s *ClientHandlerTestSuite) SetupTest() {
	s.mockUseCase = new(MockClientUseCase)
	s.handler = handlers.NewClientSSRHandler(s.mockUseCase)
}

func (s *ClientHandlerTestSuite) TestCreateClient_Success() {
	// Arrange
	formData := map[string]string{
		"name":      "John Doe",
		"cpf":       "123.456.789-00",
		"birth_day": "1990-01-01",
		"email":     "john@example.com",
		"phone":     "1234567890",
	}

	s.mockUseCase.On("CreateClient", application.CreateClientInput{
		Name:             "John Doe",
		CPF:              "123.456.789-00",
		BirthDay:         "1990-01-01",
		Email:            "john@example.com",
		Phone:            "1234567890",
		CreatorUserEmail: "creator@example.com",
	}).Return(&application.CreateClientOutput{ID: 1}, nil)

	// Act
	req := s.createFormRequest("/client", formData)
	rec := httptest.NewRecorder()
	s.handler.CreateClient(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestCreateClient_InvalidRequest() {
	// Arrange
	formData := map[string]string{
		"name": "John Doe",
		// Missing required fields
	}

	s.mockUseCase.On("CreateClient", application.CreateClientInput{
		Name:             "John Doe",
		CreatorUserEmail: "creator@example.com",
	}).Return(nil, errors.New("invalid request body"))

	// Act
	req := s.createFormRequest("/client", formData)
	rec := httptest.NewRecorder()
	s.handler.CreateClient(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestUpdateClient_Success() {
	// Arrange
	formData := map[string]string{
		"name":      "John Doe",
		"cpf":       "123.456.789-00",
		"birth_day": "1990-01-01",
		"email":     "john@example.com",
		"phone":     "1234567890",
	}

	s.mockUseCase.On("UpdateClient", application.UpdateClientInput{
		ID:               1,
		Name:             "John Doe",
		CPF:              "123.456.789-00",
		BirthDay:         "1990-01-01",
		Email:            "john@example.com",
		Phone:            "1234567890",
		CreatorUserEmail: "creator@example.com",
	}).Return(&application.UpdateClientOutput{ID: 1}, nil)

	// Act
	req := s.createFormRequest("/client/update/1", formData)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	s.handler.UpdateClient(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestUpdateClient_InvalidID() {
	// Arrange
	formData := map[string]string{
		"name": "John Doe",
	}

	// Act
	req := s.createFormRequest("/client/update/invalid", formData)
	rec := httptest.NewRecorder()
	s.handler.UpdateClient(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
}

func (s *ClientHandlerTestSuite) TestImportClients_Success() {
	// Arrange
	fileContent := []byte("name,cpf,birth_day,email,phone\nJohn Doe,123.456.789-00,1990-01-01,john@example.com,1234567890")
	fileName := "clients.csv"

	s.mockUseCase.On("ImportClients", application.ImportClientsInput{
		File:             fileContent,
		FileName:         fileName,
		CreatorUserEmail: "creator@example.com",
	}).Return(&application.ImportClientsOutput{ImportedCount: 1}, nil)

	// Act
	req := s.createMultipartRequest("/client/import", fileName, fileContent)
	rec := httptest.NewRecorder()
	s.handler.ImportClients(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestImportClients_InvalidFormat() {
	// Arrange
	fileContent := []byte("invalid,csv,format")
	fileName := "clients.txt"

	// Act
	req := s.createMultipartRequest("/client/import", fileName, fileContent)
	rec := httptest.NewRecorder()
	s.handler.ImportClients(rec, req)

	// Assert
	s.Equal(http.StatusSeeOther, rec.Code)
}

func (s *ClientHandlerTestSuite) TestListClients_Success() {
	// Arrange
	client, _ := domain.NewClient("John Doe", "123.456.789-00", "1990-01-01", "john@example.com", "1234567890")
	s.mockUseCase.On("ListClients", application.ListClientsInput{
		Term:             "John",
		Page:             1,
		PageSize:         10,
		CreatorUserEmail: "creator@example.com",
	}).Return(&application.ListClientsOutput{
		Clients:    []*domain.Client{client},
		TotalCount: 1,
		Page:       1,
		PageSize:   10,
	}, nil)

	// Act
	req := s.createQueryRequest("/client", map[string]string{
		"term":      "John",
		"page":      "1",
		"page_size": "10",
	})
	rec := httptest.NewRecorder()
	s.handler.ListClients(rec, req)

	// Assert
	s.Equal(http.StatusOK, rec.Code)
	s.mockUseCase.AssertExpectations(s.T())
}

// Helper methods
func (s *ClientHandlerTestSuite) createFormRequest(path string, formData map[string]string) *http.Request {
	form := url.Values{}
	for key, value := range formData {
		form.Add(key, value)
	}

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), utils.LoggedUserKey, "creator@example.com")
	return req.WithContext(ctx)
}

func (s *ClientHandlerTestSuite) createMultipartRequest(path string, fileName string, fileContent []byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	s.NoError(err)
	part.Write(fileContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), utils.LoggedUserKey, "creator@example.com")
	return req.WithContext(ctx)
}

func (s *ClientHandlerTestSuite) createQueryRequest(path string, queryParams map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	ctx := context.WithValue(req.Context(), utils.LoggedUserKey, "creator@example.com")
	return req.WithContext(ctx)
}

func TestClientHandlerSuite(t *testing.T) {
	suite.Run(t, new(ClientHandlerTestSuite))
}
