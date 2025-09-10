package transaction

import (
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type DbTransactionRepository interface {
	CreateTransaction(ctx context.Context, args dbchain.CreateTransactionParams, tx *sql.Tx) (dbchain.Transaction, error)
	GetListTransactionByBlockHash(ctx context.Context, bID string, tx *sql.Tx) ([]dbchain.Transaction, error)
	GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams) ([]dbchain.Transaction, error)
	CreateTxInput(ctx context.Context, args dbchain.CreateTxInputParams, tx *sql.Tx) (dbchain.TxInput, error)
	CreateTxOutput(ctx context.Context, args dbchain.CreateTxOutputParams, tx *sql.Tx) (dbchain.TxOutput, error)
	GetListTxInputByTxID(ctx context.Context, txID string) ([]dbchain.TxInput, error)
	GetListTxOutputByTxID(ctx context.Context, txID string) ([]dbchain.TxOutput, error)
	GetTxOutputByTxIDAndIndex(ctx context.Context, args dbchain.GetTxOutputByTxIDAndIndexParams, tx *sql.Tx) (dbchain.TxOutput, error)
	GetCountTransaction(ctx context.Context) (int64, error)
	FindListTxInputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxInput, error)
	FindListTxOutputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxOutput, error)
}
