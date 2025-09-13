package usecase

import (
	"context"

	"github.com/turbo514/shortenurl-v2/link/entity"
)

type ShortLinkRepository interface {
	FindByCode(ctx context.Context, code string) (*entity.ShortLink, error)
}

type EventPublisher interface {
	PublishClickEvent(ctx context.Context, event entity.ClickEvent) error
}

type ResolveShortLinkUseCase struct {
	repo      ShortLinkRepository
	publisher EventPublisher
}

func (uc *ResolveShortLinkUseCase) Resolve(ctx context.Context, code string) {

}
