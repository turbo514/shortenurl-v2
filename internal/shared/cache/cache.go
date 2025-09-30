package cache

// 虽然我觉得是没有redis的替代

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
}
