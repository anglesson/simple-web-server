package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserRepositoryTestSuite é uma suíte de testes para o GormUserRepositoryImpl.
type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.UserRepository
	tx   *gorm.DB
}

// SetupSuite é executado uma vez antes de todos os testes na suíte.
// Ele configura um banco de dados SQLite em memória.
func (s *UserRepositoryTestSuite) SetupSuite() {
	// Usando banco de dados em memória para garantir isolamento total
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(s.T(), err, "Falha ao conectar ao banco de dados em memória")

	s.db = db
	// Roda as migrações para criar a tabela 'users'
	err = s.db.AutoMigrate(&models.User{})
	require.NoError(s.T(), err, "Falha ao executar migrações")
}

// TearDownSuite é executado uma vez após todos os testes na suíte.
func (s *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

// SetupTest é executado antes de cada teste.
// Inicia uma transação e cria uma nova instância do repositório para o teste.
func (s *UserRepositoryTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = repository.NewGormUserRepository(s.tx)
}

// TearDownTest é executado após cada teste.
// Faz o rollback da transação para limpar os dados do teste.
func (s *UserRepositoryTestSuite) TearDownTest() {
	s.tx.Rollback()
}

// TestCreateUser testa a criação bem-sucedida de um usuário.
func (s *UserRepositoryTestSuite) TestCreateUser() {
	uniqueEmail := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	user, err := domain.NewUser("Any name", uniqueEmail, "Valid_Password123")
	require.NoError(s.T(), err)

	err = s.repo.Create(user)
	require.NoError(s.T(), err, "A criação do usuário não deve retornar erro")

	assert.NotZero(s.T(), user.ID, "O ID do usuário deve ser preenchido após a criação")

	// Opcional: verificar se o usuário foi realmente salvo no banco (dentro da transação)
	var foundUser models.User
	err = s.tx.First(&foundUser, user.ID).Error
	require.NoError(s.T(), err, "O usuário deve ser encontrado no banco de dados")
	assert.Equal(s.T(), uniqueEmail, foundUser.Email)
}

// TestUserRepository é a função que o Go executa, que por sua vez roda a suíte.
func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
