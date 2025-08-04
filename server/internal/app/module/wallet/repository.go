package wallet

import (
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
)

type RPCWalletRepository interface {
	GetBalance(address string) ([]byte, error)
}

type DBWalletRepository interface {
	CreateWallet(ctx context.Context, args dbwallet.CreateWalletParams, tx *sql.Tx) (dbwallet.Wallet, error)
	GetWalletByPubkey(ctx context.Context, pubkey string) (dbwallet.Wallet, error)
	IncreaseWalletBalance(ctx context.Context, args dbwallet.IncreaseWalletBalanceParams, tx *sql.Tx) error
	DecreaseWalletBalance(ctx context.Context, args dbwallet.DecreaseWalletBalanceParams, tx *sql.Tx) error
	CreateWalletAccessLog(ctx context.Context, args dbwallet.CreateWalletAccessLogParams, tx *sql.Tx) error
	GetListAccessLogByWalletID(ctx context.Context, args dbwallet.GetListAccessLogByWalletIDParams) ([]dbwallet.WalletAccessLog, error)
	GetWalletByPubKeyHash(ctx context.Context, PubkeyHash string) (dbwallet.Wallet, error)
}
