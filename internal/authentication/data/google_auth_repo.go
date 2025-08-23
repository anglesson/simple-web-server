package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/anglesson/simple-web-server/internal/authentication/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleAuthRepository é a camada de dados que lida com a API do Google.
type GoogleAuthRepository struct {
	config *oauth2.Config
}

// GoogleUserInfo representa as informações do usuário retornadas pela API do Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email,omitempty"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewGoogleAuthRepository cria uma nova instância do repositório.
func NewGoogleAuthRepository() *GoogleAuthRepository {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleAuthRepository{
		config: config,
	}
}

// GetAuthURL retorna a URL de autorização do Google
func (r *GoogleAuthRepository) GetAuthURL(state string) string {
	return r.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForUserInfo troca o código de autorização por informações do usuário usando a API real do Google.
func (r *GoogleAuthRepository) ExchangeCodeForUserInfo(code string) (*session.User, error) {
	fmt.Printf("Trocando código de autorização '%s' com a API do Google...\n", code)

	// Valida se as variáveis de ambiente estão configuradas
	if r.config.ClientID == "" || r.config.ClientSecret == "" {
		return nil, fmt.Errorf("Google OAuth não configurado: CLIENT_ID e CLIENT_SECRET são obrigatórios")
	}

	ctx := context.Background()

	// Troca o código de autorização por um token de acesso
	token, err := r.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("erro ao trocar código por token: %v", err)
	}

	// Opção 1: Usando a biblioteca oficial do Google
	userInfo, err := r.getUserInfoWithGoogleAPI(ctx, token)
	if err != nil {
		// Fallback para requisição HTTP direta
		fmt.Printf("Tentativa com API oficial falhou, usando requisição HTTP direta: %v\n", err)
		userInfo, err = r.getUserInfoWithHTTP(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter informações do usuário: %v", err)
		}
	}

	// Converte as informações do Google para o formato interno
	user := &session.User{
		ID:    fmt.Sprintf("google-%s", userInfo.ID),
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}

	fmt.Printf("Informações do usuário obtidas do Google para o e-mail: %s\n", user.Email)
	return user, nil
}

// getUserInfoWithGoogleAPI obtém informações do usuário usando a biblioteca oficial do Google
func (r *GoogleAuthRepository) getUserInfoWithGoogleAPI(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	oauth2Service, err := oauth2api.NewService(ctx, option.WithTokenSource(r.config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar serviço OAuth2: %v", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações do usuário: %v", err)
	}

	return &GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: *userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}, nil
}

// getUserInfoWithHTTP obtém informações do usuário fazendo uma requisição HTTP direta
func (r *GoogleAuthRepository) getUserInfoWithHTTP(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := r.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("erro na requisição para API do Google: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta da API do Google: status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do Google: %v", err)
	}

	return &userInfo, nil
}
