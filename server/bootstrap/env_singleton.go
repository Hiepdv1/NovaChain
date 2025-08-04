package bootstrap

import (
	"ChainServer/internal/common/env"
	"sync"
)

var (
	once         sync.Once
	appEnvConfig *env.Env
)

func AppEnv() *env.Env {
	once.Do(func() {
		appEnvConfig = env.New()
	})

	return appEnvConfig
}
