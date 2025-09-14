package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/tenant/dao/model"
	"github.com/turbo514/shortenurl-v2/tenant/dao/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

var _ ITenantService = (*TenantService)(nil)

type ITenantService interface {
	CreateTenant(ctx context.Context, name string) (*model.Tenant, error)
	CreateUser(ctx context.Context, tenantId uuid.UUID, apiKey string, name string, password string) (*model.User, error)
	Login(ctx context.Context, name string, password string) (*model.User, error)
}

type TenantService struct {
	tenantRepo repository.ITenantRepo
	userRepo   repository.IUserRepo
}

func NewTenantService(tenantRepo repository.ITenantRepo, userRepo repository.IUserRepo) *TenantService {
	return &TenantService{tenantRepo: tenantRepo, userRepo: userRepo}
}

func (s *TenantService) CreateTenant(ctx context.Context, name string) (*model.Tenant, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	tenant := &model.Tenant{
		ID:        id[:],
		Name:      name,
		ApiKey:    uuid.New().String(),
		CreatedAt: t,
		UpdatedAt: t,
	}
	if err = s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *TenantService) CreateUser(ctx context.Context, tenantId uuid.UUID, apiKey string, name string, password string) (*model.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantRepo.FindById(ctx, tenantId[:])
	if err != nil {
		return nil, fmt.Errorf("查找租户失败: %w", err)
	}

	if tenant.ApiKey != apiKey {
		return nil, fmt.Errorf("apiKey错误: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	user := &model.User{
		ID:        id[:],
		TenantID:  tenantId[:],
		Name:      name,
		Password:  string(hashPassword),
		CreatedAt: t,
		UpdatedAt: t,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *TenantService) Login(ctx context.Context, name string, password string) (*model.User, error) {
	user, err := s.userRepo.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return user, nil
}
