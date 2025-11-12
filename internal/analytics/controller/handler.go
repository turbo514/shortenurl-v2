package controller

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/analytics/service"
	analyticspb "github.com/turbo514/shortenurl-v2/shared/gen/proto/analytics"
)

type ServiceHandler struct {
	analyticspb.UnimplementedAnalyticsServiceServer
	svc *service.AnalyticsService
}

var _ analyticspb.AnalyticsServiceServer = (*ServiceHandler)(nil)

func NewServiceHandler(svc *service.AnalyticsService) *ServiceHandler {
	return &ServiceHandler{
		svc: svc,
	}
}

func (s *ServiceHandler) GetTopN(ctx context.Context, req *analyticspb.GetTopNRequest) (*analyticspb.GetTopNResponse, error) {
	res, _, err := s.svc.GetTopLinks(ctx, req.N, req.TenantId)
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetTopLinks err: %w", err)
	}
	topLinks := make([]*analyticspb.TopLink, len(res.List))
	for i := range res.List {
		topLinks[i] = &analyticspb.TopLink{
			Id:          res.List[i].ID.String(),
			OriginalUrl: res.List[i].OriginalURL,
			ClickTimes:  res.List[i].ClickTimes,
		}
	}
	return &analyticspb.GetTopNResponse{
		TopLinks: topLinks,
	}, nil
}
