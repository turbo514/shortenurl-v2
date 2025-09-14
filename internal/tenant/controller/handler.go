package controller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"github.com/turbo514/shortenurl-v2/tenant/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceHandler struct {
	tenantpb.UnimplementedTenantServiceServer
	tenantService service.ITenantService
	tokenService  service.ITokenService
}

var _ tenantpb.TenantServiceServer = (*ServiceHandler)(nil)

func NewHandler(tenantService service.ITenantService, tokenService service.ITokenService) *ServiceHandler {
	return &ServiceHandler{tenantService: tenantService, tokenService: tokenService}
}

func (h *ServiceHandler) CreateTenant(ctx context.Context, req *tenantpb.CreateTenantRequest) (*tenantpb.CreateTenantResponse, error) {
	tenant, err := h.tenantService.CreateTenant(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(tenant.ID)
	if err != nil {
		return nil, err
	}

	return &tenantpb.CreateTenantResponse{
		TenantId: id.String(),
		ApiKey:   tenant.ApiKey,
	}, nil
}
func (h *ServiceHandler) CreateUser(ctx context.Context, req *tenantpb.CreateUserRequest) (*emptypb.Empty, error) {
	tenantId, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, fmt.Errorf("uuid解析错误: %w", err)
	}
	_, err = h.tenantService.CreateUser(ctx, tenantId, req.ApiKey, req.Name, req.Password)
	if err != nil {
		return nil, fmt.Errorf("用户创建失败: %w", err)
	}

	return &emptypb.Empty{}, nil
}
func (h *ServiceHandler) Login(ctx context.Context, req *tenantpb.LoginRequest) (*tenantpb.LoginResponse, error) {
	user, err := h.tenantService.Login(ctx, req.Name, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := h.tokenService.GenerateToken(ctx, uuid.UUID(user.ID), uuid.UUID(user.TenantID))
	if err != nil {
		return nil, err
	}

	return &tenantpb.LoginResponse{
		Token: token,
	}, nil
}
