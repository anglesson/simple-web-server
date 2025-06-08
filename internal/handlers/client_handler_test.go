package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/handlers"
	"github.com/anglesson/simple-web-server/internal/infrastructure"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var _ application.ClientServicePort = (*MockClientService)(nil)

type MockFlashMessage struct {
	mock.Mock
}

var _ infrastructure.FlashMessagePort = (*MockFlashMessage)(nil)

func (m *MockFlashMessage) Success(message string) {
	m.Called(message)
}

func (m *MockFlashMessage) Error(message string) {
	m.Called(message)
}

type MockClientService struct {
	mock.Mock
}

func NewMockClientService() *MockClientService {
	return &MockClientService{}
}

func (m *MockClientService) CreateClient(input application.CreateClientInput) (*models.Client, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) FindCreatorsClientByID(clientID uint, creatorID uint) (*models.Client, error) {
	args := m.Called(clientID, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) Update(input application.UpdateClientInput) (*models.Client, error) {
	args := m.Called(input)
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) CreateBatchClient(clients []*models.Client) error {
	args := m.Called(clients)
	return args.Error(1)
}

type ClientHandlerTestSuite struct {
	suite.Suite
	sut               *handlers.ClientHandler
	mockClientService *MockClientService
	mockFlashMessage  *MockFlashMessage
	flashFactory      infrastructure.FlashMessageFactory
}

func (suite *ClientHandlerTestSuite) SetupTest() {
	suite.mockClientService = NewMockClientService()
	suite.mockFlashMessage = new(MockFlashMessage)

	suite.flashFactory = func(w http.ResponseWriter, r *http.Request) infrastructure.FlashMessagePort {
		return suite.mockFlashMessage
	}

	suite.sut = handlers.NewClientHandler(suite.mockClientService, suite.flashFactory)
}

func (suite *ClientHandlerTestSuite) TestUserNotFoundInContext() {
	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birth_date=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middlewares.UserEmailKey, nil)
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

	expectedInput := application.CreateClientInput{
		Email:        "client@mail",
		Name:         "Any Name",
		Phone:        "Any Phone",
		BirthDate:    "2004-01-01",
		EmailCreator: creatorEmail,
	}

	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birth_date=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middlewares.UserEmailKey, creatorEmail)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockClientService.On("CreateClient", expectedInput).Return(
		(*models.Client)(nil), errors.New("failed to create client due to service error")).Once()
	suite.mockFlashMessage.On("Error", "failed to create client due to service error").Return().Once()

	suite.sut.ClientCreateSubmit(rr, req)

	assert.Equal(suite.T(), http.StatusSeeOther, rr.Code)

	suite.mockClientService.AssertExpectations(suite.T())
	suite.mockFlashMessage.AssertExpectations(suite.T())
}

func (suite *ClientHandlerTestSuite) TestShouldCreateClient() {
	creatorEmail := "creator@mail"

	expectedInput := application.CreateClientInput{
		Email:        "client@mail",
		Name:         "Any Name",
		Phone:        "Any Phone",
		BirthDate:    "2004-01-01",
		EmailCreator: creatorEmail,
	}

	formData := strings.NewReader("email=client@mail&name=Any Name&phone=Any Phone&birth_date=2004-01-01")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), middlewares.UserEmailKey, creatorEmail)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	suite.mockClientService.On("CreateClient", expectedInput).Return(&models.Client{}, nil).Once()
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

	expectedInput := application.UpdateClientInput{
		CPF:          "Updated CPF",
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

	ctx := context.WithValue(req.Context(), middlewares.UserEmailKey, creatorEmail)
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
