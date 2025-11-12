package service

import (
	"context"
	"github.com/turbo514/shortenurl-v2/analytics/cqrs/query"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
	"github.com/turbo514/shortenurl-v2/analytics/dto"
)

type AnalyticsService struct {
	// queryBus
	getTopLinksHandler *query.GetTopLinksHandler
	clickCounter       domain.IClickCounter
}

func NewAnalyticsService(getTopLinksHandler *query.GetTopLinksHandler, clickCunter domain.IClickCounter) *AnalyticsService {
	return &AnalyticsService{
		getTopLinksHandler: getTopLinksHandler,
		clickCounter:       clickCunter,
	}
}

func (svc *AnalyticsService) GetTopLinks(ctx context.Context, num int64, tenant string) (*dto.TopLinks, int64, error) {
	linkToClickTimes, total, err := svc.clickCounter.GetTopToday(ctx, num, tenant)
	if err != nil {
		return nil, 0, err
	}
	lst, err := svc.getTopLinksHandler.Handle(ctx, &query.GetTopLinksQuery{LinkToClickTimes: linkToClickTimes})
	if err != nil {
		return nil, 0, err
	}
	return lst, total, nil
}
