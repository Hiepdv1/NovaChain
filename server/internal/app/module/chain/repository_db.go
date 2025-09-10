package chain

import (
	"ChainServer/internal/common/env"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type dbChainRepository struct {
	env     *env.Env
	queries *dbchain.Queries
}

func NewDBChainRepository() DBChainRepository {
	return &dbChainRepository{
		env:     env.Cfg,
		queries: dbchain.New(db.Psql),
	}
}

func (r *dbChainRepository) CreateBlock(ctx context.Context, args dbchain.CreateBlockParams, tx *sql.Tx) (dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateBlock(ctx, args)
}

func (r *dbChainRepository) GetBlockByHeight(ctx context.Context, height int64, tx *sql.Tx) (dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetBlockByHeight(ctx, height)
}

func (r *dbChainRepository) GetLastBlock(ctx context.Context, tx *sql.Tx) (dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetLastBlock(ctx)
}

func (r *dbChainRepository) GetBlockByHash(ctx context.Context, hash string, tx *sql.Tx) (dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetBlockByBID(ctx, hash)
}

func (r *dbChainRepository) GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams, tx *sql.Tx) ([]dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListBlocks(ctx, args)
}

func (r *dbChainRepository) DeleteBlockByHash(ctx context.Context, hash string, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.DeleteBlockByBID(ctx, hash)
}

func (r *dbChainRepository) ExistingBlock(ctx context.Context, hash string, tx *sql.Tx) (bool, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.IsExistingBlock(ctx, hash)
}

// func (r *dbChainRepository) FindUTXO()
