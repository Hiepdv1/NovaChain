package utxo

import (
	"ChainServer/internal/db"
	dbutxo "ChainServer/internal/db/utxo"
	"context"
	"database/sql"
)

type dbUTXORepository struct {
	queries *dbutxo.Queries
}

func NewDbUTXORepository() DbUTXORepository {
	return &dbUTXORepository{
		queries: dbutxo.New(db.Psql),
	}
}

func (r *dbUTXORepository) CreateUTXO(ctx context.Context, args dbutxo.CreateUTXOParams, tx *sql.Tx) (dbutxo.Utxo, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateUTXO(ctx, args)
}

func (r *dbUTXORepository) DeleteUTXO(ctx context.Context, args dbutxo.DeleteUTXOParams, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.DeleteUTXO(ctx, args)
}

func (r *dbUTXORepository) DeleteUTXOByBlockID(ctx context.Context, b_id string, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.DeleteUTXOsByBlock(ctx, b_id)
}

func (r *dbUTXORepository) FindUTXOs(ctx context.Context, pubKeyHash string, tx *sql.Tx) ([]dbutxo.Utxo, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetUTXOsByPubKeyHash(ctx, pubKeyHash)
}
