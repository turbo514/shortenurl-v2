package rate_limiter

import "context"

type IRateLimiter interface {
	Allow(ctx context.Context, key string, permits int64) (bool, error)
}
