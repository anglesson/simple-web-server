package client_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	client_application "github.com/anglesson/simple-web-server/internal/client/application"
	common_infrastructure "github.com/anglesson/simple-web-server/internal/common/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientUseCase struct {
	mock.Mock
}

func (m *MockClientUseCase) CreateClient(input client_application.CreateClientInput) (*client_application.CreateClientOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*client_application.CreateClientOutput), args.Error(1)
}

func TestCreateClientSSR_Success(t *testing.T) {
	expected := client_application.CreateClientInput{
		Name:             "AnyName",
		CPF:              "AnyCPF",
		BirthDay:         "AnyDate",
		Email:            "AnyEmail",
		Phone:            "AnyPhone",
		CreatorUserEmail: "any_user@mail.com",
	}
	mockUsecase := new(MockClientUseCase)

	mockUsecase.On("CreateClient", expected).Return(&client_application.CreateClientOutput{}, nil).Once()

	handler := NewSSRHandler(mockUsecase)

	formData := strings.NewReader("name=AnyName&cpf=AnyCPF&birth_day=AnyDate&email=AnyEmail&phone=AnyPhone")
	req := httptest.NewRequest(http.MethodPost, "/client", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), common_infrastructure.LoggedUserKey, "any_user@mail.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler.CreateClientSubmit(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code, "Expected status code 303 See Other for success add")
	assert.Equal(t, "/client", rr.Header().Get("Location"), "Expected redirect to root path")
	mockUsecase.AssertExpectations(t)
}
