package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TrialRegressionTestSuite struct {
	suite.Suite
	userRepository repository.UserRepository
}

func TestTrialRegressionTestSuite(t *testing.T) {
	suite.Run(t, new(TrialRegressionTestSuite))
}

func (suite *TrialRegressionTestSuite) SetupSuite() {
	database.Connect()
	suite.userRepository = repository.NewGormUserRepository(database.DB)
}

func (suite *TrialRegressionTestSuite) SetupTest() {
	// Limpar dados de teste
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM users")
}

// TestTrialMiddleware_Regression_UserWithActiveTrial_ShouldAllowAccess
// Este teste simula exatamente o cenário que estava falhando:
// - Usuário criado com subscription
// - Trial ativo (IsTrialActive = true, TrialEndDate no futuro)
// - Subscription status = "inactive" (normal para trial)
// - Usuário deve ter acesso às rotas protegidas
func (suite *TrialRegressionTestSuite) TestTrialMiddleware_Regression_UserWithActiveTrial_ShouldAllowAccess() {
	// 1. Criar usuário (simulando registro)
	user := models.NewUser("testuser", "password123", "test@example.com")
	err := suite.userRepository.Create(user)
	suite.Require().NoError(err)

	// 2. Criar subscription (simulando criação automática no CreatorService)
	subscription := models.NewSubscription(user.ID, "default_plan")
	// Garantir que o trial está ativo
	subscription.IsTrialActive = true
	subscription.TrialStartDate = time.Now()
	subscription.TrialEndDate = time.Now().AddDate(0, 0, 7) // 7 dias
	subscription.SubscriptionStatus = "inactive"            // Status normal para trial
	subscription.Origin = "web"

	err = database.DB.Create(subscription).Error
	suite.Require().NoError(err)

	// 3. Verificar que a subscription foi criada corretamente no banco
	var dbSubscription models.Subscription
	err = database.DB.Where("user_id = ?", user.ID).First(&dbSubscription).Error
	suite.Require().NoError(err)

	assert.True(suite.T(), dbSubscription.IsTrialActive)
	assert.Equal(suite.T(), "inactive", dbSubscription.SubscriptionStatus)
	assert.True(suite.T(), time.Now().Before(dbSubscription.TrialEndDate))

	// 4. Buscar usuário via repository (simulando o que o middleware faz)
	foundUser := suite.userRepository.FindByUserEmail(user.Email)
	suite.Require().NotNil(foundUser)
	suite.Require().NotNil(foundUser.Subscription)

	// 5. Verificar que os métodos do trial funcionam corretamente
	assert.True(suite.T(), foundUser.IsInTrialPeriod(), "Usuário deve estar no período trial")
	assert.False(suite.T(), foundUser.IsSubscribed(), "Usuário não deve estar inscrito")
	assert.Greater(suite.T(), foundUser.DaysLeftInTrial(), 0, "Deve ter dias restantes no trial")

	// 6. Testar o middleware com o usuário real
	req := httptest.NewRequest("GET", "/dashboard", nil)
	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, user.Email)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.TrialMiddleware(nextHandler)
	handler.ServeHTTP(w, req)

	// 7. Verificar que o acesso foi permitido
	assert.Equal(suite.T(), http.StatusOK, w.Code, "Usuário com trial ativo deve ter acesso permitido")
}

// TestTrialMiddleware_Regression_UserWithoutSubscription_ShouldRedirectToSettings
// Este teste verifica que usuários sem subscription são redirecionados
func (suite *TrialRegressionTestSuite) TestTrialMiddleware_Regression_UserWithoutSubscription_ShouldRedirectToSettings() {
	// 1. Criar usuário sem subscription
	user := models.NewUser("testuser2", "password123", "test2@example.com")
	err := suite.userRepository.Create(user)
	suite.Require().NoError(err)

	// 2. Buscar usuário via repository
	foundUser := suite.userRepository.FindByUserEmail(user.Email)
	suite.Require().NotNil(foundUser)
	assert.Nil(suite.T(), foundUser.Subscription, "Usuário não deve ter subscription")

	// 3. Verificar que os métodos do trial retornam false/0
	assert.False(suite.T(), foundUser.IsInTrialPeriod())
	assert.False(suite.T(), foundUser.IsSubscribed())
	assert.Equal(suite.T(), 0, foundUser.DaysLeftInTrial())

	// 4. Testar o middleware
	req := httptest.NewRequest("GET", "/dashboard", nil)
	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, user.Email)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.TrialMiddleware(nextHandler)
	handler.ServeHTTP(w, req)

	// 5. Verificar que foi redirecionado para settings
	assert.Equal(suite.T(), http.StatusSeeOther, w.Code)
	assert.Equal(suite.T(), "/settings", w.Header().Get("Location"))
}

// TestTrialMiddleware_Regression_UserWithExpiredTrial_ShouldRedirectToSettings
// Este teste verifica que usuários com trial expirado são redirecionados
func (suite *TrialRegressionTestSuite) TestTrialMiddleware_Regression_UserWithExpiredTrial_ShouldRedirectToSettings() {
	// 1. Criar usuário
	user := models.NewUser("testuser3", "password123", "test3@example.com")
	err := suite.userRepository.Create(user)
	suite.Require().NoError(err)

	// 2. Criar subscription com trial expirado
	subscription := models.NewSubscription(user.ID, "default_plan")
	subscription.IsTrialActive = false
	subscription.TrialEndDate = time.Now().AddDate(0, 0, -1) // Trial expirado há 1 dia
	subscription.SubscriptionStatus = "inactive"

	err = database.DB.Create(subscription).Error
	suite.Require().NoError(err)

	// 3. Buscar usuário via repository
	foundUser := suite.userRepository.FindByUserEmail(user.Email)
	suite.Require().NotNil(foundUser)
	suite.Require().NotNil(foundUser.Subscription)

	// 4. Verificar que os métodos do trial retornam false/0
	assert.False(suite.T(), foundUser.IsInTrialPeriod())
	assert.False(suite.T(), foundUser.IsSubscribed())
	assert.Equal(suite.T(), 0, foundUser.DaysLeftInTrial())

	// 5. Testar o middleware
	req := httptest.NewRequest("GET", "/dashboard", nil)
	ctx := context.WithValue(req.Context(), middleware.UserEmailKey, user.Email)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.TrialMiddleware(nextHandler)
	handler.ServeHTTP(w, req)

	// 6. Verificar que foi redirecionado para settings
	assert.Equal(suite.T(), http.StatusSeeOther, w.Code)
	assert.Equal(suite.T(), "/settings", w.Header().Get("Location"))
}

// TestUserRepository_Regression_PreloadSubscription_WorksCorrectly
// Este teste verifica especificamente que o Preload está funcionando
func (suite *TrialRegressionTestSuite) TestUserRepository_Regression_PreloadSubscription_WorksCorrectly() {
	// 1. Criar usuário com subscription
	user := models.NewUser("testuser4", "password123", "test4@example.com")
	err := suite.userRepository.Create(user)
	suite.Require().NoError(err)

	subscription := models.NewSubscription(user.ID, "test_plan")
	err = database.DB.Create(subscription).Error
	suite.Require().NoError(err)

	// 2. Buscar usuário via FindByUserEmail (que usa Preload)
	foundUser := suite.userRepository.FindByUserEmail(user.Email)
	suite.Require().NotNil(foundUser)

	// 3. Verificar que a subscription foi carregada
	assert.NotNil(suite.T(), foundUser.Subscription, "Subscription deve ser carregada via Preload")
	assert.Equal(suite.T(), user.ID, foundUser.Subscription.UserID)
	assert.Equal(suite.T(), "test_plan", foundUser.Subscription.PlanID)
	assert.True(suite.T(), foundUser.Subscription.IsTrialActive)

	// 4. Verificar que os métodos do trial funcionam
	assert.True(suite.T(), foundUser.IsInTrialPeriod())
	assert.False(suite.T(), foundUser.IsSubscribed())
	assert.Greater(suite.T(), foundUser.DaysLeftInTrial(), 0)
}
