package wallet

import (
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
)

type WalletDBRepository struct {
	queries *dbwallet.Queries
}

func NewDBWalletRepository() DBWalletRepository {
	return &WalletDBRepository{
		queries: dbwallet.New(db.Psql),
	}
}

func (r *WalletDBRepository) CreateWallet(ctx context.Context, args dbwallet.CreateWalletParams, tx *sql.Tx) (dbwallet.Wallet, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateWallet(ctx, args)
}

func (r *WalletDBRepository) GetWalletByPubKeyHash(ctx context.Context, PubkeyHash string, tx *sql.Tx) (*dbwallet.Wallet, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	wallet, err := q.GetWalletByPubKeyHash(ctx, PubkeyHash)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletDBRepository) GetWalletByPubkey(ctx context.Context, pubkey []byte, tx *sql.Tx) (*dbwallet.Wallet, error) {

	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	address := utils.PubKeyToAddress(pubkey)

	args := dbwallet.GetWalletByAddrAndPubkeyParams{
		Address:   string(address),
		PublicKey: hex.EncodeToString(pubkey),
	}

	wallet, err := q.GetWalletByAddrAndPubkey(ctx, args)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletDBRepository) IncreaseWalletBalance(ctx context.Context, args dbwallet.IncreaseWalletBalanceParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.IncreaseWalletBalance(ctx, args)
}

func (r *WalletDBRepository) IncreaseWalletBalanceByPubKeyHash(ctx context.Context, args dbwallet.IncreaseWalletBalanceByPubKeyHashParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.IncreaseWalletBalanceByPubKeyHash(ctx, args)
}

func (r *WalletDBRepository) DecreaseWalletBalance(ctx context.Context, args dbwallet.DecreaseWalletBalanceParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.DecreaseWalletBalance(ctx, args)
}

func (r *WalletDBRepository) CreateWalletAccessLog(ctx context.Context, args dbwallet.CreateWalletAccessLogParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.CreateWalletAccessLog(ctx, args)
}

func (r *WalletDBRepository) GetListAccessLogByWalletID(ctx context.Context, args dbwallet.GetListAccessLogByWalletIDParams, tx *sql.Tx) ([]dbwallet.WalletAccessLog, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	wallets, err := q.GetListAccessLogByWalletID(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			return make([]dbwallet.WalletAccessLog, 0), nil
		}
		return nil, err
	}

	return wallets, nil
}

func (r *WalletDBRepository) ExistsWalletByPubKey(ctx context.Context, pubkey []byte, tx *sql.Tx) (bool, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	addr := utils.PubKeyToAddress(pubkey)

	exists, err := q.ExistsWalletByAddrAndPubkey(ctx, dbwallet.ExistsWalletByAddrAndPubkeyParams{
		Address:   string(addr),
		PublicKey: hex.EncodeToString(pubkey),
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *WalletDBRepository) UpdateWalletInfoByWalletID(ctx context.Context, args dbwallet.UpdateWalletInfoByWalletIDParams, tx *sql.Tx) (*dbwallet.Wallet, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	wallet, err := q.UpdateWalletInfoByWalletID(ctx, args)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
