package wallet

import (
	cacheRedis "ChainServer/internal/cache/redis"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"time"
)

type RPCWalletRepository interface {
	GetBalance(address string) (*Balance, error)
}

type DBWalletRepository interface {
	CreateWallet(ctx context.Context, args dbwallet.CreateWalletParams, tx *sql.Tx) (dbwallet.Wallet, error)
	GetWalletByPubkey(ctx context.Context, pubkey []byte, tx *sql.Tx) (*dbwallet.Wallet, error)
	IncreaseWalletBalance(ctx context.Context, args dbwallet.IncreaseWalletBalanceParams, tx *sql.Tx) error
	DecreaseWalletBalance(ctx context.Context, args dbwallet.DecreaseWalletBalanceParams, tx *sql.Tx) error
	CreateWalletAccessLog(ctx context.Context, args dbwallet.CreateWalletAccessLogParams, tx *sql.Tx) error
	GetListAccessLogByWalletID(ctx context.Context, args dbwallet.GetListAccessLogByWalletIDParams, tx *sql.Tx) ([]dbwallet.WalletAccessLog, error)
	GetWalletByPubKeyHash(ctx context.Context, PubkeyHash string, tx *sql.Tx) (*dbwallet.Wallet, error)
	ExistsWalletByPubKey(ctx context.Context, pubkey []byte, tx *sql.Tx) (bool, error)
	IncreaseWalletBalanceByPubKeyHash(ctx context.Context, args dbwallet.IncreaseWalletBalanceByPubKeyHashParams, tx *sql.Tx) error
	UpdateWalletInfoByWalletID(ctx context.Context, args dbwallet.UpdateWalletInfoByWalletIDParams, tx *sql.Tx) (*dbwallet.Wallet, error)
}

type CacheWalletRepository interface {
	NewKey(parts ...string) *cacheRedis.CacheKey
	GetWalletById(ctx context.Context, walletId string) (*dbwallet.Wallet, error)
	AddWallet(ctx context.Context, wallet dbwallet.Wallet, ttl ...time.Duration) error
}
