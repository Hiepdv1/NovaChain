package chain

import (
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type RPCChainRepository interface {
	GetBlocks(startHash string, limit int) ([]*Block, error)
	GetBlocksByHeightRange(height, limit int64) ([]*Block, error)
	GetBlockByHash(hash string) (*Block, error)
	GetCommonBlock(locator [][]byte) (*Block, error)
}

type DBChainRepository interface {
	CountMiners(ctx context.Context) (int64, error)
	GetMiners(ctx context.Context, arg dbchain.GetMinersParams) ([]dbchain.GetMinersRow, error)
	DeleteBlockByRangeHeight(ctx context.Context, start, end int64, tx *sql.Tx) error
	GetBlockLocator(ctx context.Context, tx *sql.Tx) ([][]byte, error)
	GetRecentBlocksForNetworkInfo(ctx context.Context, limit int32) ([]dbchain.GetRecentBlocksForNetworkInfoRow, error)
	CountFuzzy(ctx context.Context, content string) (int64, error)
	CountFuzzyByType(ctx context.Context, content string) ([]dbchain.CountFuzzyByTypeRow, error)
	SearchExact(ctx context.Context, content string) ([]dbchain.SearchExactRow, error)
	SearchFuzzy(ctx context.Context, arg dbchain.SearchFuzzyParams) ([]dbchain.SearchFuzzyRow, error)
	CountDistinctMiners(ctx context.Context, tx *sql.Tx) (int64, error)
	GetCountTodayWorkerMiners(ctx context.Context, tx *sql.Tx) (int64, error)
	GetListBlockByHours(ctx context.Context, hours int64, tx *sql.Tx) ([]dbchain.Block, error)
	GetBestHeight(ctx context.Context, tx *sql.Tx) (int64, error)
	GetBlockCountByHours(ctx context.Context, hours int64, tx *sql.Tx) (int64, error)
	CreateBlock(ctx context.Context, args dbchain.CreateBlockParams, tx *sql.Tx) (dbchain.Block, error)
	GetBlockByHeight(ctx context.Context, height int64, tx *sql.Tx) (dbchain.Block, error)
	GetLastBlock(ctx context.Context, tx *sql.Tx) (dbchain.Block, error)
	GetBlockByHash(ctx context.Context, hash string, tx *sql.Tx) (dbchain.Block, error)
	DeleteBlockByHash(ctx context.Context, hash string, tx *sql.Tx) error
	GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams, tx *sql.Tx) ([]dbchain.GetListBlocksRow, error)
	ExistingBlock(ctx context.Context, hash string, tx *sql.Tx) (bool, error)
	GetBlockDetailWithTransactions(ctx context.Context, arg dbchain.GetBlockDetailWithTransactionsParams) (dbchain.GetBlockDetailWithTransactionsRow, error)
}
