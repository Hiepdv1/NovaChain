package dashboard

import (
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type ChainRepository interface {
	GetBestHeight(ctx context.Context, tx *sql.Tx) (int64, error)
	GetBlockCountByHours(ctx context.Context, hours int64, tx *sql.Tx) (int64, error)
	GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams, tx *sql.Tx) ([]dbchain.Block, error)
	GetBlockByHeight(ctx context.Context, height int64, tx *sql.Tx) (dbchain.Block, error)
	GetListBlockByHours(ctx context.Context, hours int64, tx *sql.Tx) ([]dbchain.Block, error)
	CountDistinctMiners(ctx context.Context, tx *sql.Tx) (int64, error)
	GetCountTodayWorkerMiners(ctx context.Context, tx *sql.Tx) (int64, error)
}

type TXRepository interface {
	CountTodayTransaction(ctx context.Context, tx *sql.Tx) (int64, error)
	CountTransactions(ctx context.Context, tx *sql.Tx) (int64, error)
	GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams, tx *sql.Tx) ([]dbchain.Transaction, error)
	GetListFullTransaction(ctx context.Context, arg dbchain.GetListFullTransactionParams, tx *sql.Tx) ([]dbchain.GetListFullTransactionRow, error)

	// Transaction Pending
	CountPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
	CountTodayPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
}
