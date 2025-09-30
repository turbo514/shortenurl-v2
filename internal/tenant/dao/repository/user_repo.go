package repository

import (
	"context"
	"github.com/turbo514/shortenurl-v2/tenant/dao/model"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Create(context.Context, *model.User) error
	FindByNameAndTenantID(ctx context.Context, name string, tenantID []byte) (*model.User, error)
}

var _ IUserRepo = (*UserRepo)(nil)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(ctx context.Context, user *model.User) error {
	return repo.db.WithContext(ctx).Create(user).Error
}

func (repo *UserRepo) FindByNameAndTenantID(ctx context.Context, name string, tenantID []byte) (*model.User, error) {
	user := &model.User{}
	if err := repo.db.WithContext(ctx).Where("tenant_id", tenantID).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
