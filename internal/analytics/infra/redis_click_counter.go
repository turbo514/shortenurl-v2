package infra

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
	"github.com/turbo514/shortenurl-v2/shared/keys"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	"github.com/turbo514/shortenurl-v2/shared/util"
	"time"
)

var _ domain.IClickCounter = (*RedisClickCounter)(nil)

type RedisClickCounter struct {
	conn *redis.Client
}

func NewRedisClickCounter(conn *redis.Client) *RedisClickCounter {
	return &RedisClickCounter{conn: conn}
}

//func (r *RedisClickCounter) Increase(ctx context.Context, linkID uuid.UUID, increasement int) error {
//	//TODO implement me
//	panic("implement me")
//}

func (r *RedisClickCounter) IncreaseMany(ctx context.Context, links []*domain.ClickEvent) error {
	pipe := r.conn.Pipeline()

	m := map[string]map[string]int{}
	for i := 0; i < len(links); i++ {
		if _, ok := m[util.UUIDToRedisMember(links[i].TenantID)]; !ok {
			m[util.UUIDToRedisMember(links[i].TenantID)] = map[string]int{}
		}
		m[util.UUIDToRedisMember(links[i].TenantID)][util.UUIDToRedisMember(links[i].LinkID)]++
	}

	for tenant, idToIncrement := range m {
		key := getRankKey(tenant)
		for id, increment := range idToIncrement {
			if err := pipe.ZIncrBy(ctx, key, float64(increment), id).Err(); err != nil {
				// FIXME
				return fmt.Errorf("pipe.ZIncrBy err: %w", err)
			}
			mylog.GetLogger().Debug("链接点击量增加", "租户", tenant, "链接", id, "增加量", increment, "Key", key)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipe.Exec err: %w", err)
	}

	return nil
}

func (r *RedisClickCounter) GetTopToday(ctx context.Context, num int64, tenant string) (map[uuid.UUID]int64, int64, error) {
	key := getRankKey(tenant)

	if num <= 0 {
		num = 100
	}

	// 获取总数
	total, err := r.conn.ZCard(ctx, key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("ZCard err: %w", err)
	}
	if total == 0 {
		return map[uuid.UUID]int64{}, 0, nil
	}

	res, err := r.conn.ZRevRangeWithScores(ctx, key, 0, num-1).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("r.conn.ZRevRangeWithScores err: %w", err)
	}

	linkToClickTimes := make(map[uuid.UUID]int64)
	for i := range res {
		id, err := util.RedisMemberToUUID(res[i].Member.(string))
		if err != nil {
			return nil, 0, fmt.Errorf("util.RedisMemberToUUID err: %w", err)
		}
		linkToClickTimes[id] = int64(res[i].Score)
	}

	mylog.GetLogger().Debug("获取到排行榜", "租户", tenant, "key", key, "排行榜", linkToClickTimes)

	return linkToClickTimes, total, nil
}

func getRankKey(tenantID string) string {
	return keys.HotLinksKey + ":" + tenantID + ":" + time.Now().Format("2006-01-02")
}
