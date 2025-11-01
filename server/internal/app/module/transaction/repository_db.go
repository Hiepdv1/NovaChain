package transaction

import (
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	"context"
	"database/sql"
)

type dbTransactionRepository struct {
	queries        *dbchain.Queries
	pendingQueries *dbPendingTx.Queries
}

func NewDbTransactionRepository() DbTransactionRepository {
	return &dbTransactionRepository{
		queries:        dbchain.New(db.Psql),
		pendingQueries: dbPendingTx.New(db.Psql),
	}
}

func (r *dbTransactionRepository) CreateTransaction(ctx context.Context, args dbchain.CreateTransactionParams, tx *sql.Tx) (dbchain.Transaction, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CreateTransaction(ctx, args)
}

func (r *dbTransactionRepository) GetListTransactionByBlockHash(ctx context.Context, arg dbchain.GetListTransactionByBIDParams, tx *sql.Tx) ([]dbchain.Transaction, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListTransactionByBID(ctx, arg)
}

func (r *dbTransactionRepository) GetFullTransactionByBlockHash(ctx context.Context, b_hash string, tx *sql.Tx) ([]dbchain.Transaction, error) {
	q := r.queries
	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.GetFullTransactionByBID(ctx, b_hash)
}

func (r *dbTransactionRepository) CountTransactionByBID(ctx context.Context, b_id string, tx *sql.Tx) (int64, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CountTransactionByBID(ctx, b_id)
}
func (r *dbTransactionRepository) GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams, tx *sql.Tx) ([]dbchain.Transaction, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListTransactions(ctx, args)
}

func (r *dbTransactionRepository) GetListFullTransaction(ctx context.Context, arg dbchain.GetListFullTransactionParams, tx *sql.Tx) ([]dbchain.GetListFullTransactionRow, error) {
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.GetListFullTransaction(ctx, arg)
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
	q := r.queries

	if tx != nil {
		q = r.queries.WithTx(tx)
	}
	return q.GetTxOutputByTxIDAndIndex(ctx, args)
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

func (r *dbTransactionRepository) CountTodayTransaction(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.queries
	if tx != nil {
		q = r.queries.WithTx(tx)
	}

	return q.CountTodayTransactions(ctx)
}

func (r *dbTransactionRepository) SearchFuzzyTransactionsByBlock(ctx context.Context, arg dbchain.SearchFuzzyTransactionsByBlockParams) ([]dbchain.Transaction, error) {
	return r.queries.SearchFuzzyTransactionsByBlock(ctx, arg)
}

func (r *dbTransactionRepository) CountFuzzyTransactionsByBlock(ctx context.Context, arg dbchain.CountFuzzyTransactionsByBlockParams) (int64, error) {
	return r.queries.CountFuzzyTransactionsByBlock(ctx, arg)
}

// ---------------- Pending Transactions ----------------

func (r *dbTransactionRepository) FindTxPendingByTxID(ctx context.Context, txID string, tx *sql.Tx) (dbPendingTx.PendingTransaction, error) {

	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.SelectTxPendingByTxID(ctx, txID)
}

func (r *dbTransactionRepository) InsertPendingTransaction(ctx context.Context, args dbPendingTx.InsertPendingTransactionParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}
	return q.InsertPendingTransaction(ctx, args)
}

func (r *dbTransactionRepository) UpdatePendingTransactionStatus(ctx context.Context, args dbPendingTx.UpdatePendingTransactionStatusParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}
	return q.UpdatePendingTransactionStatus(ctx, args)
}

func (r *dbTransactionRepository) UpdatePendingTransactionPriority(ctx context.Context, args dbPendingTx.UpdatePendingTransactionPriorityParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}
	return q.UpdatePendingTransactionPriority(ctx, args)
}

func (r *dbTransactionRepository) SelectPendingTransactions(ctx context.Context, args dbPendingTx.SelectPendingTransactionsParams, tx *sql.Tx) ([]dbPendingTx.SelectPendingTransactionsRow, error) {
	q := r.pendingQueries

	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.SelectPendingTransactions(ctx, args)
}

func (r *dbTransactionRepository) ExistingPendingTransaction(ctx context.Context, arg dbPendingTx.PendingTxExistsParams, tx *sql.Tx) (bool, error) {
	q := r.pendingQueries

	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.PendingTxExists(ctx, arg)
}

func (r *dbTransactionRepository) CountPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.pendingQueries

	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.CountPendingTxs(ctx)
}

func (r *dbTransactionRepository) CountTodayPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error) {
	q := r.pendingQueries

	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.CountTodayPendingTxs(ctx)
}

func (r *dbTransactionRepository) GetPendingTransactions(ctx context.Context, arg dbPendingTx.GetListPendingTxsParams) ([]dbPendingTx.GetListPendingTxsRow, error) {

	return r.pendingQueries.GetListPendingTxs(ctx, arg)
}

func (r *dbTransactionRepository) GetCountPendingTransaction(ctx context.Context) (int64, error) {
	return r.pendingQueries.GetCountPendingTxs(ctx)
}

// ---------------- Pending Tx Data ----------------

func (r *dbTransactionRepository) InsertPendingTxData(ctx context.Context, args dbPendingTx.InsertPendingTxDataParams, tx *sql.Tx) (dbPendingTx.PendingTxDatum, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}
	return q.InsertPendingTxData(ctx, args)
}

func (r *dbTransactionRepository) FindTxPendingByAddr(ctx context.Context, args dbPendingTx.PendingTxsByAddressParams, tx *sql.Tx) ([]dbPendingTx.PendingTxsByAddressRow, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.PendingTxsByAddress(ctx, args)
}

func (r *dbTransactionRepository) FindTxPendingByAddrAndStatus(ctx context.Context, args dbPendingTx.PendingTxsByAddressAndStatusParams, tx *sql.Tx) ([]dbPendingTx.PendingTxsByAddressAndStatusRow, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.PendingTxsByAddressAndStatus(ctx, args)
}

func (r *dbTransactionRepository) CountPendingTxByAddr(ctx context.Context, address string, tx *sql.Tx) (int64, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}
	return q.CountPendingTransactionsByAddr(ctx, address)
}

func (r *dbTransactionRepository) UpdatePendingTxsStatus(ctx context.Context, args dbPendingTx.UpdatePendingTxsStatusParams, tx *sql.Tx) (int64, error) {
	q := r.pendingQueries
	if tx != nil {
		q = r.pendingQueries.WithTx(tx)
	}

	return q.UpdatePendingTxsStatus(ctx, args)
}
