package redis

import "fmt"

type CacheNamespace string

const (
	NamespaceBlacklist CacheNamespace = "blacklist"
	NamespaceSession   CacheNamespace = "session"
	NamespaceCache     CacheNamespace = "cache"
	NamespaceRateLimit CacheNamespace = "ratelimit"
	NamespaceGeneric   CacheNamespace = "generic"
)

type CacheKey struct {
	Namespace CacheNamespace
	Key       string
}

func (ck CacheKey) String() string {
	return fmt.Sprint(string(ck.Namespace), ":", ck.Key)
}
