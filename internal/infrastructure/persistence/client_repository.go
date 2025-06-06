package persistence

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/application"
	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
	"github.com/anglesson/simple-web-server/internal/domain"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{
		db: GetDB(),
	}
}

func (r *ClientRepository) Create(client *domain.Client) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) Update(client *domain.Client) error {
	return r.db.Save(client).Error
}

func (r *ClientRepository) FindByCPF(cpf common_domain.CPF) *domain.Client {
	var client domain.Client
	if err := r.db.Where("cpf = ?", cpf).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return &client
}

func (r *ClientRepository) FindByID(id uint) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) FindByCreatorID(creatorID uint, query application.ClientQuery) ([]*domain.Client, error) {
	var clients []*domain.Client
	db := r.db.Model(&domain.Client{})

	if query.Term != "" {
		db = db.Where("name LIKE ? OR cpf LIKE ?", "%"+query.Term+"%", "%"+query.Term+"%")
	}

	if query.Pagination != nil {
		offset := (query.Pagination.Page - 1) * query.Pagination.PageSize
		db = db.Offset(offset).Limit(query.Pagination.PageSize)
	}

	if err := db.Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) CreateBatch(clients []*domain.Client) error {
	return r.db.Create(clients).Error
}

func (r *ClientRepository) List(query application.ClientQuery) ([]*domain.Client, int, error) {
	var clients []*domain.Client
	var total int64
	db := r.db.Model(&domain.Client{})

	if query.Term != "" {
		db = db.Where("name LIKE ? OR cpf LIKE ?", "%"+query.Term+"%", "%"+query.Term+"%")
	}

	if query.Pagination != nil {
		offset := (query.Pagination.Page - 1) * query.Pagination.PageSize
		db = db.Offset(offset).Limit(query.Pagination.PageSize)
	}

	if err := db.Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&domain.Client{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return clients, int(total), nil
}
