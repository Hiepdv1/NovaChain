package ratelimiter

import (
	"ChainServer/internal/cache/redis"
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key redis.CacheKey) (bool, Result, error)
	SetConfig(ctx context.Context, key redis.CacheKey, cfg Config) error
}

type Config struct {
	Rate   int
	Burst  int
	Window time.Duration
}

type Result struct {
	Remaining  int
	RetryAfter time.Duration
}
