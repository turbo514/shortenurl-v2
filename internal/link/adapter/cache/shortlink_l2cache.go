package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/link/domain"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"time"
)

// 我不感觉会有替代品,就不搞接口了
type ShortLinkL2Cache struct {
	client *redis.Client
}

func NewShortLinkL2Cache(client *redis.Client) *ShortLinkL2Cache {
	return &ShortLinkL2Cache{
		client: client,
	}
}

func (c *ShortLinkL2Cache) GetLinkByCode(ctx context.Context, code string) (*domain.ShortLink, error) {
	key := "su:links:code:" + code
	v, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("从Redis获取值[code: %s]失败: %w", code, err)
	}

	shortLink := domain.ShortLink{}
	if err := json.Unmarshal([]byte(v), &shortLink); err != nil {
		return nil, fmt.Errorf("解析Json失败: %w", err)
	}

	if shortLink.ID == uuid.Nil {
		return nil, zerr.ErrNotFoundDB
	}
	return &shortLink, nil
}

func (c *ShortLinkL2Cache) PutLinkByCode(ctx context.Context, shortlink *domain.ShortLink, ttl time.Duration) error {
	key := "su:links:code:" + shortlink.ShortCode
	v, err := json.Marshal(shortlink)
	if err != nil {
		return fmt.Errorf("json序列化失败: %w", err)
	}
	if err := c.client.Set(ctx, key, v, ttl).Err(); err != nil {
		return fmt.Errorf("向Redis写值失败: %w", err)
	}
	return nil
}
