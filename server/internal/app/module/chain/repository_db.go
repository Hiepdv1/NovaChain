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

func (r *dbChainRepository) GetBlockByHeight(ctx context.Context, height int64) (dbchain.Block, error) {
	return r.queries.GetBlockByHeight(ctx, height)
}

func (r *dbChainRepository) GetLastBlock(ctx context.Context) (dbchain.Block, error) {
	return r.queries.GetLastBlock(ctx)
}

func (r *dbChainRepository) GetBlockByHash(ctx context.Context, hash string) (dbchain.Block, error) {
	return r.queries.GetBlockByBID(ctx, hash)
}

func (r *dbChainRepository) GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams) ([]dbchain.Block, error) {
	return r.queries.GetListBlocks(ctx, args)
}

func (r *dbChainRepository) DeleteBlockByHash(ctx context.Context, hash string, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.DeleteBlockByBID(ctx, hash)
}
