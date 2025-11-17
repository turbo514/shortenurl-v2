package controller

import (
	"context"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
)

type IAcker interface {
	Ack(ctx context.Context) error
	Nack(ctx context.Context) error
}

type ClickEventWithAcker struct {
	Event []*domain.ClickEvent
	Acker IAcker
}
