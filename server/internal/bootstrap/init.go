package bootstrap

import (
	"ChainServer/internal/cache/redis"
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/gobtypes"
	"ChainServer/internal/common/logger"
	"time"
)

func Init() {
	env.InitEnv()

	logger.InitLogger(env.Cfg.AppEnv)

	redis.InitRedis(config.RedisConfig{
		URL:          "redis://localhost:6379",
		MaxRetries:   10,
		RetryBackoff: 3 * time.Second,
		Namespaces: []config.RedisNamespaceConfig{
			{Name: string(redis.NamespaceBlacklist), TTL: time.Hour + 2*time.Minute},
			{Name: string(redis.NamespaceSession), TTL: time.Hour + 2*time.Minute},
			{Name: string(redis.NamespaceCache), TTL: 24 * time.Hour},
			{Name: string(redis.NamespaceRateLimit), TTL: 1 * time.Minute},
			{Name: string(redis.NamespaceGeneric), TTL: 12 * time.Hour},
		},
	})

	gobtypes.Init()
}
