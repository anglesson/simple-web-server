package repository

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(file *models.File) error
	FindByID(id uint) (*models.File, error)
	FindByCreator(creatorID uint) ([]*models.File, error)
	Update(file *models.File) error
	Delete(id uint) error
	FindByType(creatorID uint, fileType string) ([]*models.File, error)
	FindActiveByCreator(creatorID uint) ([]*models.File, error)
}

type GormFileRepository struct {
	db *gorm.DB
}

func NewGormFileRepository(db *gorm.DB) *GormFileRepository {
	return &GormFileRepository{db: db}
}

func (r *GormFileRepository) Create(file *models.File) error {
	log.Printf("Criando arquivo: Nome=%s, CreatorID=%d, Tipo=%s", file.Name, file.CreatorID, file.FileType)
	err := r.db.Create(file).Error
	log.Printf("Arquivo criado com sucesso, ID=%d, erro: %v", file.Model.ID, err)
	return err
}

func (r *GormFileRepository) FindByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.Preload("Creator").First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *GormFileRepository) FindByCreator(creatorID uint) ([]*models.File, error) {
	var files []*models.File
	log.Printf("Executando consulta FindByCreator para creatorID: %d", creatorID)
	err := r.db.Where("creator_id = ?", creatorID).Order("created_at DESC").Find(&files).Error
	log.Printf("Consulta FindByCreator retornou %d arquivos, erro: %v", len(files), err)
	return files, err
}

func (r *GormFileRepository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

func (r *GormFileRepository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

func (r *GormFileRepository) FindByType(creatorID uint, fileType string) ([]*models.File, error) {
	var files []*models.File
	err := r.db.Where("creator_id = ? AND file_type = ?", creatorID, fileType).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (r *GormFileRepository) FindActiveByCreator(creatorID uint) ([]*models.File, error) {
	var files []*models.File
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, true).Order("created_at DESC").Find(&files).Error
	return files, err
}
