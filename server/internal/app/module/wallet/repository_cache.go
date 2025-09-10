package wallet

import (
	cacheRedis "ChainServer/internal/cache/redis"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const NamespaceWallet cacheRedis.CacheNamespace = "wallet"

type walletCacheRepository struct {
	rdb       *redis.Client
	namespace cacheRedis.CacheNamespace
	ttl       time.Duration
}

func NewWalletCacheRepository(rdb *redis.Client) CacheWalletRepository {
	return &walletCacheRepository{
		rdb:       rdb,
		namespace: NamespaceWallet,
		ttl:       10 * time.Minute,
	}
}

func (w *walletCacheRepository) NewKey(parts ...string) *cacheRedis.CacheKey {
	return &cacheRedis.CacheKey{
		Namespace: w.namespace,
		Key:       strings.Join(parts, ":"),
	}
}

func (w *walletCacheRepository) GetWalletById(ctx context.Context, pubkeyHex string) (*dbwallet.Wallet, error) {
	key := w.NewKey(pubkeyHex, "auth")

	wallet, err := cacheRedis.GetTyped[dbwallet.Wallet](ctx, *key)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (w *walletCacheRepository) AddWallet(ctx context.Context, wallet dbwallet.Wallet, ttl ...time.Duration) error {
	key := w.NewKey(wallet.PublicKey, "auth")

	effectiveTTL := w.ttl

	if len(ttl) > 0 && ttl[0] > 0 {
		effectiveTTL = ttl[0]
	}

	return cacheRedis.Set(ctx, *key, wallet, effectiveTTL)
}
