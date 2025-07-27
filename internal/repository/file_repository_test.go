package repository_test

import (
	"testing"
	"time"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type FileRepositoryTestSuite struct {
	suite.Suite
	db             *gorm.DB
	fileRepository repository.FileRepository
	creator        *models.Creator
}

func (suite *FileRepositoryTestSuite) SetupSuite() {
	// Configurar banco de teste
	database.Connect()
	suite.db = database.DB

	// Auto-migrate
	suite.db.AutoMigrate(&models.File{}, &models.Creator{}, &models.User{})
}

func (suite *FileRepositoryTestSuite) SetupTest() {
	// Limpar dados antes de cada teste
	suite.db.Exec("DELETE FROM files")
	suite.db.Exec("DELETE FROM creators")
	suite.db.Exec("DELETE FROM users")

	// Criar criador de teste
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	suite.db.Create(user)

	birthDate, _ := time.Parse("2006-01-02", "1990-01-01")
	suite.creator = &models.Creator{
		Name:      "Test Creator",
		Email:     "creator@example.com",
		CPF:       "12345678901",
		Phone:     "11999999999",
		BirthDate: birthDate,
		UserID:    user.ID,
	}
	suite.db.Create(suite.creator)

	// Inicializar repositório
	suite.fileRepository = repository.NewGormFileRepository(suite.db)
}

func (suite *FileRepositoryTestSuite) TearDownSuite() {
	// Limpar após todos os testes
	suite.db.Exec("DELETE FROM files")
	suite.db.Exec("DELETE FROM creators")
	suite.db.Exec("DELETE FROM users")
}

func (suite *FileRepositoryTestSuite) TestCreate() {
	// Arrange
	file := &models.File{
		Name:         "test-file.pdf",
		OriginalName: "original-test.pdf",
		Description:  "Test file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/test-file.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/test-file.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}

	// Act
	err := suite.fileRepository.Create(file)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), file.ID)

	// Verificar se foi salvo no banco
	var savedFile models.File
	suite.db.First(&savedFile, file.ID)
	assert.Equal(suite.T(), file.Name, savedFile.Name)
}

func (suite *FileRepositoryTestSuite) TestFindByID() {
	// Arrange
	file := &models.File{
		Name:         "test-file.pdf",
		OriginalName: "original-test.pdf",
		Description:  "Test file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/test-file.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/test-file.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(file)

	// Act
	foundFile, err := suite.fileRepository.FindByID(file.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundFile)
	assert.Equal(suite.T(), file.Name, foundFile.Name)
	assert.Equal(suite.T(), file.CreatorID, foundFile.CreatorID)
}

func (suite *FileRepositoryTestSuite) TestFindByCreator() {
	// Arrange
	file1 := &models.File{
		Name:         "file1.pdf",
		OriginalName: "original1.pdf",
		Description:  "First file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/file1.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/file1.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	file2 := &models.File{
		Name:         "file2.pdf",
		OriginalName: "original2.pdf",
		Description:  "Second file",
		FileType:     "pdf",
		FileSize:     2048 * 1024,
		S3Key:        "files/1/file2.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/file2.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(file1)
	suite.db.Create(file2)

	// Act
	files, err := suite.fileRepository.FindByCreator(suite.creator.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), files, 2)
	assert.Equal(suite.T(), file2.Name, files[0].Name) // Ordenado por created_at DESC
	assert.Equal(suite.T(), file1.Name, files[1].Name)
}

func (suite *FileRepositoryTestSuite) TestFindByType() {
	// Arrange
	pdfFile := &models.File{
		Name:         "file.pdf",
		OriginalName: "original.pdf",
		Description:  "PDF file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/file.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/file.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	imageFile := &models.File{
		Name:         "image.jpg",
		OriginalName: "original.jpg",
		Description:  "Image file",
		FileType:     "image",
		FileSize:     512 * 1024,
		S3Key:        "files/1/image.jpg",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/image.jpg",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(pdfFile)
	suite.db.Create(imageFile)

	// Act
	pdfFiles, err := suite.fileRepository.FindByType(suite.creator.ID, "pdf")

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), pdfFiles, 1)
	assert.Equal(suite.T(), "pdf", pdfFiles[0].FileType)
}

func (suite *FileRepositoryTestSuite) TestFindActiveByCreator() {
	// Arrange
	activeFile := &models.File{
		Name:         "active.pdf",
		OriginalName: "original.pdf",
		Description:  "Active file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/active.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/active.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	inactiveFile := &models.File{
		Name:         "inactive.pdf",
		OriginalName: "original.pdf",
		Description:  "Inactive file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/inactive.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/inactive.pdf",
		Status:       false,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(activeFile)
	suite.db.Create(inactiveFile)

	// Act
	activeFiles, err := suite.fileRepository.FindActiveByCreator(suite.creator.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), activeFiles, 1)
	assert.Equal(suite.T(), activeFile.Name, activeFiles[0].Name)
	assert.True(suite.T(), activeFiles[0].Status)
}

func (suite *FileRepositoryTestSuite) TestUpdate() {
	// Arrange
	file := &models.File{
		Name:         "test-file.pdf",
		OriginalName: "original-test.pdf",
		Description:  "Original description",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/test-file.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/test-file.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(file)

	// Act
	file.Description = "Updated description"
	err := suite.fileRepository.Update(file)

	// Assert
	assert.NoError(suite.T(), err)

	// Verificar se foi atualizado no banco
	var updatedFile models.File
	suite.db.First(&updatedFile, file.ID)
	assert.Equal(suite.T(), "Updated description", updatedFile.Description)
}

func (suite *FileRepositoryTestSuite) TestDelete() {
	// Arrange
	file := &models.File{
		Name:         "test-file.pdf",
		OriginalName: "original-test.pdf",
		Description:  "Test file",
		FileType:     "pdf",
		FileSize:     1024 * 1024,
		S3Key:        "files/1/test-file.pdf",
		S3URL:        "https://bucket.s3.amazonaws.com/files/1/test-file.pdf",
		Status:       true,
		CreatorID:    suite.creator.ID,
	}
	suite.db.Create(file)

	// Act
	err := suite.fileRepository.Delete(file.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verificar se foi deletado do banco
	var deletedFile models.File
	result := suite.db.First(&deletedFile, file.ID)
	assert.Error(suite.T(), result.Error) // Deve retornar erro pois não existe mais
}

func TestFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileRepositoryTestSuite))
}
