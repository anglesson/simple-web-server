package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks_service "github.com/anglesson/simple-web-server/internal/service/mocks"
	mocks_cookies "github.com/anglesson/simple-web-server/pkg/cookie/mocks"

	handler "github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/service"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// MockTemplateRenderer implements template.TemplateRenderer for testing
type MockTemplateRenderer struct {
	mock.Mock
}

func (m *MockTemplateRenderer) View(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}, layout string) {
	m.Called(w, r, page, data, layout)
}

func (m *MockTemplateRenderer) ViewWithoutLayout(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	m.Called(w, r, page, data)
}

var _ service.ClientService = (*mocks_service.MockClientService)(nil)
var _ template.TemplateRenderer = (*MockTemplateRenderer)(nil)

type ClientHandlerTestSuite struct {
	suite.Suite
	sut                  *handler.ClientHandler
	mockClientService    *mocks_service.MockClientService
	mockCreatorService   *mocks_service.MockCreatorService
	mockFlashMessage     *mocks_cookies.MockFlashMessage
	mockTemplateRenderer *MockTemplateRenderer
	flashFactory         web.FlashMessageFactory
}

func (suite *ClientHandlerTestSuite) SetupTest() {
	suite.mockClientService = mocks_service.NewMockClientService()
	suite.mockFlashMessage = new(mocks_cookies.MockFlashMessage)
	suite.mockCreatorService = new(mocks_service.MockCreatorService)
	suite.mockTemplateRenderer = new(MockTemplateRenderer)

	suite.flashFactory = func(w http.ResponseWriter, r *http.Request) web.FlashMessagePort {
		return suite.mockFlashMessage
	}

	suite.sut = handler.NewClientHandler(suite.mockClientService, suite.mockCreatorService, suite.flashFactory, suite.mockTemplateRenderer)
}

func (suite *ClientHandlerTestSuite) TestUserNotFoundInContext() {
	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birth_date=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, nil)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockFlashMessage.On("Error", "Unauthorized. Invalid user email").Return().Once()

	suite.mockClientService.AssertNotCalled(suite.T(), "CreateClient", mock.Anything)

	suite.sut.ClientCreateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)

	suite.mockFlashMessage.AssertExpectations(suite.T())
	suite.mockClientService.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestShouldRedirectBackIfErrorsOnService() {
	creatorEmail := "creator@mail"

	expectedInput := service.CreateClientInput{
		Email:        "client@mail",
		Name:         "Any Name",
		Phone:        "Any Phone",
		BirthDate:    "2004-01-01",
		EmailCreator: creatorEmail,
	}

	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birthdate=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, creatorEmail)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockClientService.On("CreateClient", expectedInput).Return(
		(*service.CreateClientOutput)(nil), errors.New("failed to create client due to service error")).Once()
	suite.mockFlashMessage.On("Error", "failed to create client due to service error").Return().Once()

	suite.sut.ClientCreateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusSeeOther, rr.Code)

	suite.mockClientService.AssertExpectations(suite.T())
	suite.mockFlashMessage.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestShouldCreateClient() {
	creatorEmail := "creator@mail"

	expectedInput := service.CreateClientInput{
		Email:        "client@mail",
		Name:         "Any Name",
		Phone:        "Any Phone",
		BirthDate:    "2004-01-01",
		EmailCreator: creatorEmail,
	}

	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birthdate=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, creatorEmail)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockClientService.On("CreateClient", expectedInput).Return(&service.CreateClientOutput{}, nil).Once()
	suite.mockFlashMessage.On("Success", "Cliente foi cadastrado!").Return().Once()

	suite.sut.ClientCreateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusSeeOther, rr.Code)
	assert.Equal(suite.T(), "/client", rr.Header().Get("Location"))

	suite.mockClientService.AssertExpectations(suite.T())
	suite.mockFlashMessage.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestShouldUpdateClientSuccessfully() {
	creatorEmail := "creator@mail"
	clientID := uint(1)

	expectedInput := service.UpdateClientInput{
		ID:           clientID,
		Email:        "updated@mail.com",
		Phone:        "Updated Phone",
		EmailCreator: creatorEmail,
	}

	expectedClient := &models.Client{
		Model: gorm.Model{ID: clientID},
	}

	formData := strings.NewReader("cpf=Updated CPF&email=updated@mail.com&phone=Updated Phone")
	req := httptest.NewRequest(http.MethodPost, "/client/update/1", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set chi route context for URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, creatorEmail)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockClientService.On("Update", expectedInput).Return(expectedClient, nil).Once()
	suite.mockFlashMessage.On("Success", "Cliente foi atualizado!").Return().Once()

	suite.sut.ClientUpdateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusSeeOther, rr.Code)
	assert.Equal(suite.T(), "/client", rr.Header().Get("Location"))

	suite.mockClientService.AssertExpectations(suite.T())
	suite.mockFlashMessage.AssertExpectations(suite.T())
}

func TestClientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ClientHandlerTestSuite))
}

func TestClient_GetInitials(t *testing.T) {
	tests := []struct {
		name     string
		client   models.Client
		expected string
	}{
		{
			name:     "Two names",
			client:   models.Client{Name: "João Silva"},
			expected: "JS",
		},
		{
			name:     "Single name",
			client:   models.Client{Name: "João"},
			expected: "J",
		},
		{
			name:     "Three names",
			client:   models.Client{Name: "João Pedro Silva"},
			expected: "JS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.client.GetInitials()
			assert.Equal(t, tt.expected, result)
		})
	}
}
