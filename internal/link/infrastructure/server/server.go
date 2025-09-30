package server

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/link/usecase"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	"time"
)

var _ linkpb.LinkServiceServer = (*GrpcServer)(nil)

type GrpcServer struct {
	linkpb.UnimplementedLinkServiceServer
	service *usecase.LinkUseCase
}

func NewGrpcServer(service *usecase.LinkUseCase) *GrpcServer {
	return &GrpcServer{
		service: service,
	}
}

func (s *GrpcServer) CreateLink(ctx context.Context, req *linkpb.CreateLinkRequest) (*linkpb.CreateLinkResponse, error) {
	// TODO: 参数校验
	shortlink, err := s.service.Shorten(ctx, req.OriginalUrl, req.TenantId, req.UserId, req.Expiration.AsDuration())
	if err != nil {
		return nil, fmt.Errorf("创建短链接失败: %w", err)
	}
	return &linkpb.CreateLinkResponse{
		OriginalUrl: shortlink.OriginalURL,
		ShortCode:   shortlink.ShortCode,
	}, nil
}
func (s *GrpcServer) ResolveLink(ctx context.Context, req *linkpb.ResolveLinkRequest) (*linkpb.ResolveLinkResponse, error) {
	// TODO: 参数校验
	originalUrl, err := s.service.Resolve(ctx, usecase.ResolveRequest{
		Code:      req.ShortCode,
		ClickTime: time.Now(),
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	})
	if err != nil {
		return nil, fmt.Errorf("处理短链接解析失败: %w", err)
	}

	return &linkpb.ResolveLinkResponse{
		OriginalUrl: originalUrl,
	}, nil
}
