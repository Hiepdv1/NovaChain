package wallet

import (
	"ChainServer/internal/common/utils"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
)

type dbWalletRepository struct {
	queries *dbwallet.Queries
}

func NewDBWalletRepository(db *sql.DB) DBWalletRepository {
	return &dbWalletRepository{
		queries: dbwallet.New(db),
	}
}

func (r *dbWalletRepository) CreateWallet(ctx context.Context, args dbwallet.CreateWalletParams, tx *sql.Tx) (dbwallet.Wallet, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateWallet(ctx, args)
}

func (r *dbWalletRepository) GetWalletByPubKeyHash(ctx context.Context, PubkeyHash string) (dbwallet.Wallet, error) {
	return r.queries.GetWalletByPubKeyHash(ctx, PubkeyHash)
}

func (r *dbWalletRepository) GetWalletByPubkey(ctx context.Context, pubkey string) (dbwallet.Wallet, error) {
	address := utils.PubKeyToAddress([]byte(pubkey))

	args := dbwallet.GetWalletByAddrAndPubkeyParams{
		Address:   address,
		PublicKey: pubkey,
	}

	return r.queries.GetWalletByAddrAndPubkey(ctx, args)
}

func (r *dbWalletRepository) IncreaseWalletBalance(ctx context.Context, args dbwallet.IncreaseWalletBalanceParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.IncreaseWalletBalance(ctx, args)
}

func (r *dbWalletRepository) DecreaseWalletBalance(ctx context.Context, args dbwallet.DecreaseWalletBalanceParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.DecreaseWalletBalance(ctx, args)
}

func (r *dbWalletRepository) CreateWalletAccessLog(ctx context.Context, args dbwallet.CreateWalletAccessLogParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.CreateWalletAccessLog(ctx, args)
}

func (r *dbWalletRepository) GetListAccessLogByWalletID(ctx context.Context, args dbwallet.GetListAccessLogByWalletIDParams) ([]dbwallet.WalletAccessLog, error) {
	return r.queries.GetListAccessLogByWalletID(ctx, args)
}
