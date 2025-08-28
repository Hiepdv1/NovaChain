package redis

import (
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/utils"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	Client   *redis.Client
	once     sync.Once
	redisCfg config.RedisConfig
	clientMu sync.RWMutex
)

func InitRedis(cfg config.RedisConfig) {
	once.Do(func() {
		opt, err := redis.ParseURL(cfg.URL)
		if err != nil {
			log.Error(fmt.Errorf("failed to parse redis url %v", err))
			return
		}

		opt.MaxRetries = int(cfg.MaxRetries)
		opt.MaxRetryBackoff = cfg.RetryBackoff

		clientMu.Lock()
		Client = redis.NewClient(opt)
		clientMu.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = Client.Ping(ctx).Result()
		if err != nil {
			log.Error(fmt.Errorf("failed to connect redis %v, Config: %+v", err, cfg))
			return
		}

		redisCfg = cfg
		log.Info("✅ Connected to redis successfully")
	})

}

func Close() error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if Client == nil {
		return errors.New("redis client has not been initialized yet")
	}
	return Client.Close()
}

func reconnect(ctx context.Context) error {
	log.Infof("Attempting to reconnect to redis with url: %s", redisCfg.URL)
	otp, err := redis.ParseURL(redisCfg.URL)
	if err != nil {
		return err
	}

	otp.MaxRetries = int(redisCfg.MaxRetries)
	otp.MaxRetryBackoff = redisCfg.RetryBackoff

	clientMu.Lock()
	Client = redis.NewClient(otp)
	clientMu.Unlock()

	_, err = Client.Ping(ctx).Result()
	if err != nil {
		log.Errorf("Reconnect failed: %v", err)
		return err
	}

	log.Info("Reconnected to redis successfully")
	return nil
}

func withRetry(ctx context.Context, fn func() error) error {
	err := fn()
	if err == nil {
		return nil
	}

	if err == redis.ErrClosed || err.Error() == "EOF" {
		if reconnectErr := reconnect(ctx); reconnectErr != nil {
			return reconnectErr
		}
		return fn()
	}

	return err
}

func getDefaultTTL(namespace CacheNamespace) time.Duration {

	for _, ns := range redisCfg.Namespaces {
		if ns.Name == string(namespace) {
			return ns.TTL
		}
	}

	return time.Hour
}

func Set(ctx context.Context, key CacheKey, value any, ttl ...time.Duration) error {
	data, err := utils.GobEncode(value)
	if err != nil {
		return err
	}

	effectiveTTL := getDefaultTTL(key.Namespace)
	if len(ttl) > 0 && ttl[0] > 0 {
		effectiveTTL = ttl[0]
	}

	return withRetry(ctx, func() error {
		clientMu.RLock()
		defer clientMu.RUnlock()

		err := Client.Set(ctx, key.String(), data, effectiveTTL).Err()
		if err != nil {
			log.Errorf("Failed to set key %s: %v", key.String(), err)
		}

		return err
	})
}

func GetTyped[T any](ctx context.Context, key CacheKey) (T, error) {
	var result T

	err := withRetry(ctx, func() error {
		clientMu.RLock()
		defer clientMu.RUnlock()

		data, err := Client.Get(ctx, key.String()).Bytes()
		if err != nil {
			if err == redis.Nil {
				return err
			}
			log.Errorf("Failed to get key %s: %v", key.String(), err)
			return err
		}

		result, err = utils.GobDecode[T](data)

		return err
	})

	return result, err
}

func Get(ctx context.Context, key CacheKey) (string, error) {
	var result string
	err := withRetry(ctx, func() error {
		clientMu.RLock()
		defer clientMu.RUnlock()

		data, err := Client.Get(ctx, key.String()).Bytes()
		if err != nil && err != redis.Nil {
			log.Errorf("Failed to get key %s: %v", key.String(), err)
		}

		result, err = utils.GobDecode[string](data)

		return err
	})

	return result, err
}

func Del(ctx context.Context, key CacheKey) error {
	return withRetry(ctx, func() error {
		clientMu.RLock()
		defer clientMu.RUnlock()

		err := Client.Del(ctx, key.String()).Err()
		if err != nil {
			log.Errorf("Failed to delete key %s: %v", key.String(), err)
		}

		return err
	})
}

func Exists(ctx context.Context, key CacheKey) (bool, error) {
	var exists int64

	err := withRetry(ctx, func() error {
		clientMu.RLock()
		defer clientMu.RUnlock()

		var err error

		exists, err = Client.Exists(ctx, key.String()).Result()
		if err != nil {
			log.Errorf("Failed to check existence of key %s: %v", key.String(), err)
		}

		return err
	})

	return exists > 0, err
}
