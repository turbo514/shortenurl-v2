package repository

import (
	"context"
	"github.com/turbo514/shortenurl-v2/tenant/dao/model"
	"gorm.io/gorm"
)

var _ ITenantRepo = (*TenantRepo)(nil)

type ITenantRepo interface {
	Create(context.Context, *model.Tenant) error
	FindById(context.Context, []byte) (*model.Tenant, error)
}

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (repo *TenantRepo) Create(ctx context.Context, tenant *model.Tenant) error {
	return repo.db.WithContext(ctx).Create(tenant).Error
}

func (repo *TenantRepo) FindById(ctx context.Context, id []byte) (*model.Tenant, error) {
	tenant := model.Tenant{}
	err := repo.db.WithContext(ctx).First(&tenant, "id = ?", id).Error
	return &tenant, err
}
