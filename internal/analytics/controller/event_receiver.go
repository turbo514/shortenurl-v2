package controller

import (
	"context"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
)

type IEventReceiver[EventType any] interface {
	StartAConsumer(ctx context.Context) (<-chan EventType, error)
}

type IAcker interface {
	Ack() error
	Nack() error
}

type ClickEventWithAcker struct {
	Event *domain.ClickEvent
	Acker IAcker
}
