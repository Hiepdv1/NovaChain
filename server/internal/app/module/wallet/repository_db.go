package wallet

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type dbWalletRepository struct {
	queries *dbwallet.Queries
}

func NewDBWalletRepository() DBWalletRepository {
	return &dbWalletRepository{
		queries: dbwallet.New(db.Psql),
	}
}

func (r *dbWalletRepository) CreateWallet(ctx context.Context, args dbwallet.CreateWalletParams, tx *sql.Tx) (dbwallet.Wallet, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateWallet(ctx, args)
}

func (r *dbWalletRepository) GetWalletByPubKeyHash(ctx context.Context, PubkeyHash string) (*dbwallet.Wallet, error) {
	wallet, err := r.queries.GetWalletByPubKeyHash(ctx, PubkeyHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NotFound(fmt.Sprintf("Wallet not found with pubkey hash %s", PubkeyHash), nil)
		}
		return nil, err
	}

	return &wallet, nil
}

func (r *dbWalletRepository) GetWalletByPubkey(ctx context.Context, pubkey []byte) (*dbwallet.Wallet, error) {
	address := utils.PubKeyToAddress(pubkey)

	args := dbwallet.GetWalletByAddrAndPubkeyParams{
		Address:   string(address),
		PublicKey: hex.EncodeToString(pubkey),
	}

	wallet, err := r.queries.GetWalletByAddrAndPubkey(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NotFound(fmt.Sprintf("Wallet not found with pubkey %s", hex.EncodeToString(pubkey)), nil)
		}

		return nil, err
	}

	return &wallet, nil
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
	wallets, err := r.queries.GetListAccessLogByWalletID(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			return make([]dbwallet.WalletAccessLog, 0), nil
		}
		return nil, err
	}

	return wallets, nil
}

func (r *dbWalletRepository) ExistsWalletByPubKey(ctx context.Context, pubkey []byte) (bool, error) {
	addr := utils.PubKeyToAddress(pubkey)

	exists, err := r.queries.ExistsWalletByAddrAndPubkey(ctx, dbwallet.ExistsWalletByAddrAndPubkeyParams{
		Address:   string(addr),
		PublicKey: hex.EncodeToString(pubkey),
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}
