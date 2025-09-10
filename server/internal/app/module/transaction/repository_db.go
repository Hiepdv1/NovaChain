package transaction

import (
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type dbTransactionRepository struct {
	queries *dbchain.Queries
}

func NewDbTransactionRepository() DbTransactionRepository {
	return &dbTransactionRepository{
		queries: dbchain.New(db.Psql),
	}
}

func (r *dbTransactionRepository) CreateTransaction(ctx context.Context, args dbchain.CreateTransactionParams, tx *sql.Tx) (dbchain.Transaction, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateTransaction(ctx, args)
}

func (r *dbTransactionRepository) GetListTransactionByBlockHash(ctx context.Context, bID string, tx *sql.Tx) ([]dbchain.Transaction, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListTransactionByBID(ctx, bID)
}

func (r *dbTransactionRepository) GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams) ([]dbchain.Transaction, error) {
	return r.queries.GetListTransactions(ctx, args)
}

func (r *dbTransactionRepository) CreateTxInput(ctx context.Context, args dbchain.CreateTxInputParams, tx *sql.Tx) (dbchain.TxInput, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.CreateTxInput(ctx, args)
}

func (r *dbTransactionRepository) CreateTxOutput(ctx context.Context, args dbchain.CreateTxOutputParams, tx *sql.Tx) (dbchain.TxOutput, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateTxOutput(ctx, args)
}

func (r *dbTransactionRepository) GetListTxInputByTxID(ctx context.Context, txID string) ([]dbchain.TxInput, error) {
	return r.queries.GetListTxInputByTxID(ctx, txID)
}

func (r *dbTransactionRepository) GetListTxOutputByTxID(ctx context.Context, txID string) ([]dbchain.TxOutput, error) {
	return r.queries.GetListTxOutputByTxId(ctx, txID)
}

func (r *dbTransactionRepository) GetTxOutputByTxIDAndIndex(ctx context.Context, args dbchain.GetTxOutputByTxIDAndIndexParams, tx *sql.Tx) (dbchain.TxOutput, error) {
	return r.queries.GetTxOutputByTxIDAndIndex(ctx, args)
}

func (r *dbTransactionRepository) GetCountTransaction(ctx context.Context) (int64, error) {
	return r.queries.GetCountTransaction(ctx)
}

func (r *dbTransactionRepository) FindListTxInputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxInput, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.FindTxInputByBlockID(ctx, b_id)
}

func (r *dbTransactionRepository) FindListTxOutputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxOutput, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.FindListTxOutputByBlockID(ctx, b_id)
}
