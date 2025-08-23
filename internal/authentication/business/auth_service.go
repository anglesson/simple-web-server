package business

import (
	"fmt"

	"github.com/anglesson/simple-web-server/internal/authentication/data"
	"github.com/anglesson/simple-web-server/internal/authentication/session"
)

// AuthService contém a lógica de negócio do módulo de autenticação.
type AuthService struct {
	googleAuthRepo *data.GoogleAuthRepository
	sessionStore   *session.SessionStore
}

// NewAuthService cria uma instância de AuthService com as suas dependências.
func NewAuthService(
	googleAuthRepo *data.GoogleAuthRepository,
	sessionStore *session.SessionStore,
) *AuthService {
	return &AuthService{
		googleAuthRepo: googleAuthRepo,
		sessionStore:   sessionStore,
	}
}

// HandleGoogleLogin é a função que orquestra o processo de login com o Google.
func (s *AuthService) HandleGoogleLogin(authorizationCode string) (string, error) {
	// 1. Usa a camada de dados para obter as informações do usuário do Google.
	// O serviço não se importa com a forma como isso é feito, apenas que a informação é retornada.
	user, err := s.googleAuthRepo.ExchangeCodeForUserInfo(authorizationCode)
	if err != nil {
		fmt.Printf("Erro na autenticação com o Google: %v\n", err)
		return "", fmt.Errorf("falha na autenticação")
	}

	// 2. Com as informações do usuário em mãos, usa a camada de sessão para criar uma sessão.
	sessionID, err := s.sessionStore.CreateSession(*user)
	if err != nil {
		fmt.Printf("Erro ao criar a sessão: %v\n", err)
		return "", fmt.Errorf("falha ao criar a sessão")
	}

	// 3. Retorna o ID da sessão para a camada de apresentação.
	fmt.Printf("Login bem-sucedido. ID da sessão retornado: %s\n", sessionID)
	return sessionID, nil
}
