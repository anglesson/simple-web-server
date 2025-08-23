package session

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// User representa as informações básicas do usuário que queremos armazenar na sessão.
// Isso evita que a aplicação precise fazer requisições externas a todo momento.
type User struct {
	ID    string
	Email string
	Name  string
}

// Session representa uma sessão do usuário, contendo as informações dele.
type Session struct {
	User User
}

// SessionStore gerencia o armazenamento de sessões.
type SessionStore struct {
	store map[string]Session
	mu    sync.Mutex
}

// NewSessionStore cria uma nova instância de SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		store: make(map[string]Session),
	}
}

// CreateSession cria uma nova sessão para o usuário e retorna o ID da sessão.
func (s *SessionStore) CreateSession(user User) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := uuid.New().String()
	s.store[sessionID] = Session{User: user}

	fmt.Printf("Sessão criada com sucesso para o usuário: %s\n", user.Email)
	return sessionID, nil
}

// GetSessionUser busca as informações do usuário a partir do ID da sessão.
// Este é o método que o middleware chamará.
func (s *SessionStore) GetSessionUser(sessionID string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.store[sessionID]
	if !ok {
		return nil, fmt.Errorf("sessão não encontrada para o ID: %s", sessionID)
	}

	return &session.User, nil
}

// DeleteSession remove uma sessão do armazenamento. Usada para o logout.
func (s *SessionStore) DeleteSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[sessionID]; !ok {
		return fmt.Errorf("sessão não encontrada para o ID: %s", sessionID)
	}

	delete(s.store, sessionID)
	fmt.Printf("Sessão %s removida com sucesso.\n", sessionID)
	return nil
}
