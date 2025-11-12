package adapter

import (
	"context"
	"github.com/turbo514/shortenurl-v2/link/domain"
)

type IShortLinkRepository interface {
	FindByCode(ctx context.Context, code string) (*domain.ShortLink, error)
	CreateLink(ctx context.Context, shortLink *domain.ShortLink) error
}
