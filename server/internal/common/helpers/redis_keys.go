package helpers

import (
	"ChainServer/internal/cache/redis"
	"fmt"
)

type KeyType string

const (
	AuthKeyTypeJWT    KeyType = "auth:jwt"
	AuthKeyTypeSig    KeyType = "auth:sig"
	RequestKeyTypeSig KeyType = "req:sig"
)

func BlacklistSigKey(sigType KeyType, sig string) redis.CacheKey {
	return redis.CacheKey{
		Namespace: redis.NamespaceBlacklist,
		Key:       fmt.Sprintf("%s:%s", sigType, sig),
	}
}
