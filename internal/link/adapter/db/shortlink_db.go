package db

import (
	"context"
	"github.com/turbo514/shortenurl-v2/link/entity"
)

type IShortLinkDB interface {
	FindByCode(ctx context.Context, code string) (*entity.ShortLink, error)
	CreateLink(ctx context.Context, shortLink *entity.ShortLink) error
}
