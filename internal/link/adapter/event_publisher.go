package adapter

import (
	"context"
	"github.com/turbo514/shortenurl-v2/shared/dto"
)

type IEventPublisher interface {
	PublishClickEvent(ctx context.Context, event dto.ClickEvent) error
	PublishCreateEvent(ctx context.Context, event dto.CreateEvent) error
}
