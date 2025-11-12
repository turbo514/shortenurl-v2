package rate_limiter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const luaScript = `
-- KEYS[1]: 限流器的唯一键 (e.g., "ratelimiter:user123")
-- ARGV[1]: 令牌桶容量 (capacity)
-- ARGV[2]: 令牌生成速率 (rate, 每秒生成的数量)
-- ARGV[3]: 当前时间戳 (now, 浮点数表示的秒)
-- ARGV[4]: 本次请求消耗的令牌数 (permits)

local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local permits = tonumber(ARGV[4])

local fields = redis.call('HMGET', key, 'tokens', 'last_refill_time')
local current_tokens = tonumber(fields[1])
local last_refill_time = tonumber(fields[2])

-- 如果是第一次访问，进行初始化
if current_tokens == nil then
  current_tokens = capacity
  last_refill_time = now
else
  -- 计算自上次填充以来需要补充的令牌数
  local elapsed = now - last_refill_time
  local tokens_to_add = elapsed * rate
  
  -- 填充令牌，但不能超过容量上限
  current_tokens = math.min(capacity, current_tokens + tokens_to_add)
end

local allowed = false
-- 如果当前令牌数足够本次消耗
if current_tokens >= permits then
  -- 消耗令牌
  current_tokens = current_tokens - permits
  allowed = true
end

-- 无论是否允许，都更新令牌数和最后刷新时间
-- 这是关键的修正：将 last_refill_time 总是更新为当前时间戳 now
redis.call('HSET', key, 'tokens', current_tokens, 'last_refill_time', now)

-- 为 key 设置一个合理的过期时间，防止冷数据永久占用内存
local expire_time = math.ceil(capacity / rate) * 2
redis.call('EXPIRE', key, expire_time)

if allowed then
  return 1
else
  -- 返回剩余需要等待的时间（可选，但很有用）
  -- 如果桶是空的，计算需要多久才能攒够 permits 个令牌
  local wait_time = (permits - current_tokens) / rate
  -- 返回一个负数表示拒绝，绝对值可以表示建议的等待秒数
  return -wait_time
  -- 或者简单返回 0
  -- return 0
end
`
const inf float64 = -1.0

type RedisRateLimiter struct {
	client *redis.Client
	script *redis.Script
	// 令牌桶容量
	capacity int64
	// 令牌生成速率 (每秒)
	rate float64
}

// NewRedisRateLimiter 创建一个新的限流器实例
func NewRedisRateLimiter(client *redis.Client, rate float64, capacity int64) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:   client,
		script:   redis.NewScript(luaScript),
		capacity: capacity,
		rate:     rate,
	}
}

// Allow 检查是否允许指定数量的请求
// key: 需要限流的资源的唯一标识 (例如用户ID, API路径)
// permits: 本次请求需要消耗的令牌数 (通常为1)
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, permits int64) (bool, error) {
	if r.rate == inf {
		return true, nil
	}

	// 构建存储在Redis中的完整键名
	redisKey := key
	now := float64(time.Now().UnixNano()) / 1e9 // 秒为单位，浮点数精度足够

	// 执行Lua脚本
	// EVALSHA 会优先使用脚本的SHA1摘要来执行，如果Redis没有缓存则会自动加载，效率更高
	result, err := r.script.Run(ctx, r.client, []string{redisKey}, r.capacity, r.rate, now, permits).Result()
	if err != nil {
		return false, fmt.Errorf("failed to execute redis script: %w", err)
	}

	// Lua脚本返回1表示允许，0表示拒绝
	return result.(int64) == 1, nil
}
