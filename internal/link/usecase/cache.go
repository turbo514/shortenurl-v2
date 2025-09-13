package usecase

import "context"

type Cache interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (any, error)
}
