package handler_test

import (
	"context"
	"errors"
	mocks_service "github.com/anglesson/simple-web-server/internal/services/mocks"
	mocks_cookies "github.com/anglesson/simple-web-server/pkg/cookie/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handler "github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/services"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var _ services.ClientService = (*mocks_service.MockClientService)(nil)

type ClientHandlerTestSuite struct {
	suite.Suite
	sut                *handler.ClientHandler
	mockClientService  *mocks_service.MockClientService
	mockCreatorService *mocks_service.MockCreatorService
	mockFlashMessage   *mocks_cookies.MockFlashMessage
	flashFactory       web.FlashMessageFactory
}

func (suite *ClientHandlerTestSuite) SetupTest() {
	suite.mockClientService = mocks_service.NewMockClientService()
	suite.mockFlashMessage = new(mocks_cookies.MockFlashMessage)
	suite.mockCreatorService = new(mocks_service.MockCreatorService)

	suite.flashFactory = func(w http.ResponseWriter, r *http.Request) web.FlashMessagePort {
		return suite.mockFlashMessage
	}

	suite.sut = handler.NewClientHandler(suite.mockClientService, suite.mockCreatorService, suite.flashFactory)
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

	expectedInput := services.CreateClientInput{
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
		(*services.CreateClientOutput)(nil), errors.New("failed to create client due to service error")).Once()
	suite.mockFlashMessage.On("Error", "failed to create client due to service error").Return().Once()

	suite.sut.ClientCreateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusSeeOther, rr.Code)

	suite.mockClientService.AssertExpectations(suite.T())
	suite.mockFlashMessage.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestShouldCreateClient() {
	creatorEmail := "creator@mail"

	expectedInput := services.CreateClientInput{
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

	suite.mockClientService.On("CreateClient", expectedInput).Return(&services.CreateClientOutput{}, nil).Once()
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

	expectedInput := services.UpdateClientInput{
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
