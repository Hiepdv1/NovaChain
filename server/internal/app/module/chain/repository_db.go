package chain

import (
	"ChainServer/internal/common/env"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
	"encoding/hex"
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

func (r *dbChainRepository) GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams, tx *sql.Tx) ([]dbchain.GetListBlocksRow, error) {
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

func (r *dbChainRepository) GetBestHeight(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetBestHeight(ctx)
}

func (r *dbChainRepository) GetBlockCountByHours(ctx context.Context, hours int64, tx *sql.Tx) (int64, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetBlockCountByHours(ctx, hours)
}

func (r *dbChainRepository) GetListBlockByHours(ctx context.Context, hours int64, tx *sql.Tx) ([]dbchain.Block, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListBlockByHours(ctx, hours)
}

func (r *dbChainRepository) CountDistinctMiners(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CountDistinctMiners(ctx)
}

func (r *dbChainRepository) GetCountTodayWorkerMiners(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetCountTodayWorkerMiners(ctx)
}

func (r *dbChainRepository) SearchExact(ctx context.Context, content string) ([]dbchain.SearchExactRow, error) {

	return r.queries.SearchExact(ctx, content)

}

func (r *dbChainRepository) SearchFuzzy(ctx context.Context, arg dbchain.SearchFuzzyParams) ([]dbchain.SearchFuzzyRow, error) {

	return r.queries.SearchFuzzy(ctx, arg)
}

func (r *dbChainRepository) CountFuzzy(ctx context.Context, content string) (int64, error) {
	return r.queries.CountFuzzy(ctx, content)
}

func (r *dbChainRepository) CountFuzzyByType(ctx context.Context, content string) ([]dbchain.CountFuzzyByTypeRow, error) {
	return r.queries.CountFuzzyByType(ctx, content)
}

func (r *dbChainRepository) GetBlockDetailWithTransactions(ctx context.Context, arg dbchain.GetBlockDetailWithTransactionsParams) (dbchain.GetBlockDetailWithTransactionsRow, error) {

	return r.queries.GetBlockDetailWithTransactions(ctx, arg)
}

func (r *dbChainRepository) GetRecentBlocksForNetworkInfo(ctx context.Context, limit int32) ([]dbchain.GetRecentBlocksForNetworkInfoRow, error) {
	return r.queries.GetRecentBlocksForNetworkInfo(ctx, limit)
}

func (r *dbChainRepository) GetBlockLocator(ctx context.Context, tx *sql.Tx) ([][]byte, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	var locator [][]byte
	step := int64(1)
	height, err := q.GetBestHeight(ctx)

	if err != nil {
		return nil, err
	}

	for height > 1 {
		block, err := q.GetBlockByHeight(ctx, height)
		if err != nil {
			return nil, err
		}
		hash, err := hex.DecodeString(block.BID)
		if err != nil {
			return nil, err
		}

		locator = append(locator, hash)

		if len(locator) > 10 {
			step *= 2
		}

		if height > step {
			height -= step
		} else {
			break
		}
	}

	genesis, _ := q.GetBlockByHeight(ctx, 1)
	hash, err := hex.DecodeString(genesis.BID)
	if err != nil {
		return nil, err
	}

	locator = append(locator, hash)

	return locator, nil
}

func (r *dbChainRepository) DeleteBlockByRangeHeight(ctx context.Context, start, end int64, tx *sql.Tx) error {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	for start > end {
		err := q.DeleteBlockByHeight(ctx, start)
		if err != nil {
			return err
		}

		start++
	}

	return nil
}
