package cache

import (
	"context"
	"fmt"
	"github.com/maypok86/otter/v2"
	"github.com/turbo514/shortenurl-v2/link/entity"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"time"
)

type ShortLinkL1Cache struct {
	cache *otter.Cache[string, entity.ShortLink]
}

func NewShortLinkL1Cache(options *otter.Options[string, entity.ShortLink]) (*ShortLinkL1Cache, error) {
	cache, err := otter.New(options)
	if err != nil {
		return nil, err
	}
	return &ShortLinkL1Cache{
		cache: cache,
	}, nil
}

// 需要确保返回的结构体只读
// id为uuid.Nil代表这是不存在的值
func (c *ShortLinkL1Cache) GetLinkByCode(ctx context.Context, code string) (*entity.ShortLink, error) {
	key := "su:links:code:" + code
	val, exist := c.cache.GetIfPresent(key)
	if !exist {
		return nil, fmt.Errorf("找不到短链接[code: %s]: %w", code, zerr.ErrNotFoundCache)
	}
	if val.IsInvalid() {
		return nil, zerr.ErrNotFoundDB
	}
	return &val, nil
}

func (c *ShortLinkL1Cache) PutLinkByCode(ctx context.Context, shortlink *entity.ShortLink, ttl time.Duration) bool {
	key := "su:links:code:" + shortlink.ShortCode
	_, ok := c.cache.Set(key, *shortlink)
	return ok
}
