package dashboard

import (
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type ChainRepository interface {
	GetBestHeight(ctx context.Context, tx *sql.Tx) (int64, error)
	GetBlockCountByHours(ctx context.Context, hours int64, tx *sql.Tx) (int64, error)
	GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams, tx *sql.Tx) ([]dbchain.GetListBlocksRow, error)
	GetBlockByHeight(ctx context.Context, height int64, tx *sql.Tx) (dbchain.Block, error)
	GetListBlockByHours(ctx context.Context, hours int64, tx *sql.Tx) ([]dbchain.Block, error)
	CountDistinctMiners(ctx context.Context, tx *sql.Tx) (int64, error)
	GetCountTodayWorkerMiners(ctx context.Context, tx *sql.Tx) (int64, error)
	GetRecentBlocksForNetworkInfo(ctx context.Context, limit int32) ([]dbchain.GetRecentBlocksForNetworkInfoRow, error)
}

type TXRepository interface {
	CountTodayTransaction(ctx context.Context, tx *sql.Tx) (int64, error)
	GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams, tx *sql.Tx) ([]dbchain.Transaction, error)
	GetListFullTransaction(ctx context.Context, arg dbchain.GetListFullTransactionParams, tx *sql.Tx) ([]dbchain.GetListFullTransactionRow, error)
	GetCountTransaction(ctx context.Context) (int64, error)

	// Transaction Pending
	CountPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
	CountTodayPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
}
