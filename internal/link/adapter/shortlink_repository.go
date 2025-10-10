package adapter

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/link/adapter/cache"
	"github.com/turbo514/shortenurl-v2/link/adapter/db"
	"github.com/turbo514/shortenurl-v2/link/entity"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "ShortLinkRepository.FindByCode")
	defer span.End()

	span.SetAttributes(attribute.String("link.short_code", code))

	ctx, dbspan := tr.Start(ctx, "")
	defer dbspan.End()

	shortlink, err := s.db.FindByCode(ctx, code)
	if err != nil {
		dbspan.SetStatus(codes.Error, "Repository查询短链失败")
		return nil, fmt.Errorf("查询短链失败: %w", err)
	}
	return shortlink, nil
}

func (s ShortLinkRepository) CreateLink(ctx context.Context, shortLink *entity.ShortLink) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "ShortLinkRepository.CreateLink")
	defer span.End()

	if err := s.db.CreateLink(ctx, shortLink); err != nil {
		span.SetStatus(codes.Error, "创建短链接失败")
		return fmt.Errorf("创建短链接失败: %w", err)
	}
	return nil
}
