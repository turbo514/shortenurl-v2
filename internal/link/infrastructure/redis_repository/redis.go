package redis_repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/domain"
	"github.com/turbo514/shortenurl-v2/link/metrics"
	"github.com/turbo514/shortenurl-v2/shared/keys"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"time"
)

var _ adapter.IShortLinkRepository = (*RedisCacheRepository)(nil)

type RedisCacheRepository struct {
	client *redis.Client
	next   adapter.IShortLinkRepository
}

func NewRedisCacheRepository(client *redis.Client, next adapter.IShortLinkRepository) *RedisCacheRepository {
	return &RedisCacheRepository{client: client, next: next}
}

func (r *RedisCacheRepository) FindByCode(ctx context.Context, code string) (*domain.ShortLink, error) {
	ctx, span := mytrace.GetTracer().Start(ctx, "RedisCacheRepository.FindByCode")
	defer span.End()

	key := keys.LinkCacheKey + ":" + code
	span.SetAttributes(
		semconv.DBSystemNameRedis,
		//semconv.DBOperationName("get"),
		attribute.String("cache.key", key),
	)

	// 尝试从Redis根据code获取对应短链
	// 如果存在且非空,则正常返回
	// 如果存在但为空,则返回不存在
	// 如果存在但解析失败,或者不存在,则继续
	start := time.Now()
	if v, err := r.client.Get(ctx, key).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			span.RecordError(err)
		}
	} else {
		var shortLink domain.ShortLink
		if err := json.Unmarshal([]byte(v), &shortLink); err == nil {
			// 检查缓存的是否空值,或者已过期
			if !shortLink.IsInvalid() && !shortLink.IsExpired() {
				metrics.AddDistributedCacheLookUpTotalHits()
				return &shortLink, nil
			} else {
				// 如果是,则直接报告未找到
				mylog.GetLogger().Debug("从redis找到空值对象")
				metrics.AddDistributedCacheLookUpTotalHitNil()
				return nil, zerr.ErrNotExist
			}
		}
	}
	end := time.Now()
	metrics.ObserveDistributedCacheDurationSecondsGet(end.Sub(start))

	// miss
	metrics.AddDistributedCacheLookUpTotalMiss()
	shortlink, err := r.next.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, zerr.ErrNotExist) {
			// 数据库中不存在该值,将空值缓存到redis中
			mylog.GetLogger().Debug("短链不存在,缓存至redis")
			v, err := json.Marshal(domain.ShortLink{ID: uuid.Nil})
			if err != nil {
				span.RecordError(err)
				mylog.GetLogger().Warn("redis FindByCode marshal shortlink failed", "err", err.Error())
				return nil, zerr.ErrNotExist
			}
			if err := r.client.Set(ctx, key, v, time.Minute).Err(); err != nil {
				span.RecordError(err)
				mylog.GetLogger().Debug("缓存空链接到redis失败", "err", err.Error())
				return nil, zerr.ErrNotExist
			}
			return nil, zerr.ErrNotExist
		} else {
			return nil, err
		}
	}

	//序列化并缓存到redis中
	v, err := json.Marshal(shortlink)
	if err != nil {
		span.RecordError(err)
		mylog.GetLogger().Warn("redis FindByCode marshal shortlink failed", "err", err.Error())
		return shortlink, nil
	}

	start = time.Now()
	if err := r.client.Set(ctx, key, v, time.Minute).Err(); err != nil {
		span.RecordError(err)
		mylog.GetLogger().Warn("redis FindByCode set shortlink failed", "err", err.Error())
		return shortlink, nil
	}
	end = time.Now()
	metrics.ObserveDistributedCacheDurationSecondsSet(end.Sub(start))

	// 返回
	return shortlink, nil
}

func (r *RedisCacheRepository) CreateLink(ctx context.Context, shortLink *domain.ShortLink) error {
	ctx, span := mytrace.GetTracer().Start(ctx, "RedisCacheRepository.CreateLink")
	defer span.End()

	if err := r.next.CreateLink(ctx, shortLink); err != nil {
		return err
	}

	// 序列化并存到redis中
	// FIXME: 这里考虑到短链是几乎不可能被更改的,所以没有采用先写数据库,再删缓存的方式
	v, err := json.Marshal(shortLink)
	if err != nil {
		span.RecordError(err)
		mylog.GetLogger().Warn("redis CreateLink marshal shortlink failed", "err", err.Error())
		return nil
	}

	key := keys.LinkCacheKey + ":" + shortLink.ShortCode
	start := time.Now()
	if err := r.client.Set(ctx, key, v, time.Minute).Err(); err != nil {
		span.RecordError(err)
		mylog.GetLogger().Warn("redis CreateLink set shortlink failed", "err", err.Error())
		return nil
	}
	end := time.Now()
	metrics.ObserveDistributedCacheDurationSecondsSet(end.Sub(start))

	return nil
}
