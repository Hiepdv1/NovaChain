package config

import "time"

type RedisNamespaceConfig struct {
	Name string
	TTL  time.Duration
}

type RedisConfig struct {
	URL          string
	MaxRetries   uint
	RetryBackoff time.Duration
	Namespaces   []RedisNamespaceConfig
}
