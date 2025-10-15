package ratelimiter

import (
	clientRedis "ChainServer/internal/cache/redis"
	"ChainServer/internal/common/apperror"
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	clientMu sync.RWMutex
)

type TokenBucketRateLimiter struct {
	rate   int
	burst  int
	period time.Duration
	mu     sync.Mutex
}

type BucketState struct {
	Tokens    float64
	Timestamp time.Time
}

func NewTokenBucketRateLimiter(cfg Config) (*TokenBucketRateLimiter, error) {

	if cfg.Rate <= 0 {
		return nil, apperror.BadRequest("Rate must be positive", nil)
	}

	if cfg.Burst <= 0 {
		return nil, apperror.BadRequest("Burst must be positive", nil)
	}

	clientMu.RLock()

	if clientRedis.Client == nil {
		clientMu.RUnlock()
		log.Error("Redis client not initialized")
		return nil, apperror.Internal("Something went wrong, please try again", nil)
	}

	clientMu.RUnlock()

	return &TokenBucketRateLimiter{
		rate:   cfg.Rate,
		burst:  cfg.Burst,
		period: time.Second / time.Duration(cfg.Rate),
	}, nil
}

func (rl *TokenBucketRateLimiter) Allow(ctx context.Context, key clientRedis.CacheKey) (bool, Result, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	state, err := clientRedis.GetTyped[BucketState](ctx, key)

	if err == redis.Nil {
		state = BucketState{
			Tokens:    float64(rl.burst),
			Timestamp: time.Now(),
		}
	} else if err != nil {
		log.Errorf("Token bucket: Failed to get state for key %s: %v", key.String(), err)
		return false, Result{}, apperror.Internal("Something went wrong, please try again.", nil)
	}

	now := time.Now()
	elapsed := now.Sub(state.Timestamp)
	newTokens := float64(elapsed) / float64(rl.period)
	state.Tokens = math.Min(float64(rl.burst), state.Tokens+newTokens)
	state.Timestamp = now

	if state.Tokens < 1 {
		retryAfter := rl.period - elapsed
		return false, Result{Remaining: 0, RetryAfter: retryAfter}, apperror.TooManyRequests(fmt.Sprintf("Rate limit exceeded, retry after %v", retryAfter), nil)
	}

	state.Tokens--
	err = clientRedis.Set(ctx, key, state, time.Minute)
	if err != nil {
		log.Errorf("Token bucket: Failed to set state for key %s: %v", key.String(), err)
		return false, Result{}, apperror.Internal("Something went wrong, please try again.", nil)
	}

	return true, Result{Remaining: int(state.Tokens)}, nil
}

func (rl *TokenBucketRateLimiter) SetConfig(ctx context.Context, key clientRedis.CacheKey, cfg Config) error {
	if cfg.Rate <= 0 {
		return apperror.BadRequest("Rate must be positive", nil)
	}

	if cfg.Burst <= 0 {
		return apperror.BadRequest("Burst must be positive", nil)
	}

	rl.mu.Lock()
	defer rl.mu.Lock()

	rl.rate = cfg.Rate
	rl.burst = cfg.Burst
	rl.period = time.Second / time.Duration(rl.burst)

	state, err := clientRedis.GetTyped[BucketState](ctx, key)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Errorf("Token bucket: Failed to set state for key %s: %v", key.String(), err)
		return apperror.Internal("Something went wrong, please try again.", nil)
	}

	state.Tokens = math.Min(float64(rl.burst), state.Tokens)
	return clientRedis.Set(ctx, key, state, time.Minute)
}
