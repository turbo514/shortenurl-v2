package otter_repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/maypok86/otter/v2"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/domain"
	"github.com/turbo514/shortenurl-v2/link/metrics"
	"github.com/turbo514/shortenurl-v2/shared/keys"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
)

var _ adapter.IShortLinkRepository = (*OtterCacheRepository)(nil)

type OtterCacheRepository struct {
	cache *otter.Cache[string, *domain.ShortLink]
	next  adapter.IShortLinkRepository
}

func NewOtterCacheRepository(cache *otter.Cache[string, *domain.ShortLink], next adapter.IShortLinkRepository) *OtterCacheRepository {
	return &OtterCacheRepository{
		cache: cache,
		next:  next,
	}
}

func (o *OtterCacheRepository) FindByCode(ctx context.Context, code string) (*domain.ShortLink, error) {
	ctx, span := mytrace.GetTracer().Start(ctx, "OtterCacheRepository.FindByCode")
	defer span.End()
	
	key := keys.LinkCacheKey + ":" + code

	if shortlink, exist := o.cache.GetIfPresent(key); exist {
		// 缓存的是空值或者已过期
		if shortlink.IsInvalid() || shortlink.IsExpired() {
			metrics.AddLocalCacheLookUpTotalHitNil()
			return nil, zerr.ErrNotExist
		} else {
			metrics.AddLocalCacheLookUpTotalHits()
			return shortlink, nil
		}
	}

	metrics.AddLocalCacheLookUpTotalHits()
	shortlink, err := o.next.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, zerr.ErrNotExist) {
			mylog.GetLogger().Debug("短链不存在,缓存至本地缓存")
			o.cache.Set(key, &domain.ShortLink{ID: uuid.Nil})
		}
		return nil, err
	}

	o.cache.Set(key, shortlink)

	return shortlink, nil
}

func (o *OtterCacheRepository) CreateLink(ctx context.Context, shortLink *domain.ShortLink) error {
	ctx, span := mytrace.GetTracer().Start(ctx, "OtterCacheRepository.CreateLink")
	defer span.End()

	if err := o.next.CreateLink(ctx, shortLink); err != nil {
		return err
	}

	key := keys.LinkCacheKey + ":" + shortLink.ShortCode

	o.cache.Set(key, shortLink)

	return nil
}
