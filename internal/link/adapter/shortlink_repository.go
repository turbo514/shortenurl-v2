package adapter

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/link/adapter/cache"
	"github.com/turbo514/shortenurl-v2/link/adapter/db"
	"github.com/turbo514/shortenurl-v2/link/entity"
)

type ShortLinkRepository struct {
	db      db.IShortLinkDB
	l1Cache *cache.ShortLinkL1Cache
	l2Cache *cache.ShortLinkL2Cache
}

func NewShortLinkRepository(db db.IShortLinkDB, l1Cache *cache.ShortLinkL1Cache, l2Cache *cache.ShortLinkL2Cache) *ShortLinkRepository {
	return &ShortLinkRepository{
		db:      db,
		l1Cache: l1Cache,
		l2Cache: l2Cache,
	}
}

func (s ShortLinkRepository) FindByCode(ctx context.Context, code string) (*entity.ShortLink, error) {
	shortlink, err := s.db.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("该短链接不存在: %w", err)
	}
	return shortlink, nil
}

func (s ShortLinkRepository) CreateLink(ctx context.Context, shortLink *entity.ShortLink) error {
	if err := s.db.CreateLink(ctx, shortLink); err != nil {
		return fmt.Errorf("添加短链接失败: %w", err)
	}
	return nil
}
