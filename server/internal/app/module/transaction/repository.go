package transaction

import (
	"ChainServer/internal/common/dto"
	dbchain "ChainServer/internal/db/chain"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	dbutxo "ChainServer/internal/db/utxo"
	"context"
	"database/sql"
)

type RpcTransactionRepository interface {
	SendTx(txs []Transaction) (*RPCSendTxResponse, error)
	FetchMiningTxIDs() (*RPCGetMiningTxResponse[string], error)
	FetchMiningTxsFull() (*RPCGetMiningTxResponse[dto.Transaction], error)
}

type DbTransactionRepository interface {
	// ---------------- Chain Transactions ----------------
	GetDetailTransaction(ctx context.Context, tx_hash string) (dbchain.GetDetailTxRow, error)
	GetCountRecentTransaction(ctx context.Context, pub_key_hash string) (int64, error)
	GetRecentTransaction(ctx context.Context, arg dbchain.GetRecentTransactionParams) ([]dbchain.GetRecentTransactionRow, error)
	GetFullTransactionByBlockHash(ctx context.Context, b_hash string, tx *sql.Tx) ([]dbchain.Transaction, error)
	CountTransactionByBID(ctx context.Context, b_id string, tx *sql.Tx) (int64, error)
	CreateTransaction(ctx context.Context, args dbchain.CreateTransactionParams, tx *sql.Tx) (dbchain.Transaction, error)
	GetListTransactionByBlockHash(ctx context.Context, arg dbchain.GetListTransactionByBIDParams, tx *sql.Tx) ([]dbchain.Transaction, error)
	GetListTransaction(ctx context.Context, args dbchain.GetListTransactionsParams, tx *sql.Tx) ([]dbchain.Transaction, error)
	CreateTxInput(ctx context.Context, args dbchain.CreateTxInputParams, tx *sql.Tx) (dbchain.TxInput, error)
	CreateTxOutput(ctx context.Context, args dbchain.CreateTxOutputParams, tx *sql.Tx) (dbchain.TxOutput, error)
	GetListTxInputByTxID(ctx context.Context, txID string) ([]dbchain.TxInput, error)
	GetListTxOutputByTxID(ctx context.Context, txID string) ([]dbchain.TxOutput, error)
	GetTxOutputByTxIDAndIndex(ctx context.Context, args dbchain.GetTxOutputByTxIDAndIndexParams, tx *sql.Tx) (dbchain.TxOutput, error)
	GetCountTransaction(ctx context.Context) (int64, error)
	FindListTxInputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxInput, error)
	FindListTxOutputByBlockHash(ctx context.Context, b_id string, tx *sql.Tx) ([]dbchain.TxOutput, error)
	CountTodayTransaction(ctx context.Context, tx *sql.Tx) (int64, error)
	GetListFullTransaction(ctx context.Context, arg dbchain.GetListFullTransactionParams, tx *sql.Tx) ([]dbchain.GetListFullTransactionRow, error)
	SearchFuzzyTransactionsByBlock(ctx context.Context, arg dbchain.SearchFuzzyTransactionsByBlockParams) ([]dbchain.Transaction, error)
	CountFuzzyTransactionsByBlock(ctx context.Context, arg dbchain.CountFuzzyTransactionsByBlockParams) (int64, error)
	GetTxSummaryByPubkeyHash(ctx context.Context, pub_key_hash string) (dbchain.GetTxSummaryByPubKeyHashRow, error)

	// ---------------- Pending Transactions ----------------
	GetCountPendingTxsByStatus(ctx context.Context, arg []string) (int64, error)
	GetPendingTxByStatus(ctx context.Context, arg dbPendingTx.GetPendingTxsByStatusParams) ([]dbPendingTx.GetPendingTxsByStatusRow, error)
	CountPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
	CountTodayPendingTxs(ctx context.Context, tx *sql.Tx) (int64, error)
	InsertPendingTransaction(ctx context.Context, args dbPendingTx.InsertPendingTransactionParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error)
	UpdatePendingTransactionStatus(ctx context.Context, args dbPendingTx.UpdatePendingTransactionStatusParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error)
	UpdatePendingTransactionPriority(ctx context.Context, args dbPendingTx.UpdatePendingTransactionPriorityParams, tx *sql.Tx) (dbPendingTx.PendingTransaction, error)
	SelectPendingTransactions(ctx context.Context, args dbPendingTx.SelectPendingTransactionsParams, tx *sql.Tx) ([]dbPendingTx.SelectPendingTransactionsRow, error)
	ExistingPendingTransaction(ctx context.Context, arg dbPendingTx.PendingTxExistsParams, tx *sql.Tx) (bool, error)
	FindTxPendingByTxID(ctx context.Context, txID string, tx *sql.Tx) (dbPendingTx.PendingTransaction, error)
	FindTxPendingByAddrAndStatus(ctx context.Context, args dbPendingTx.PendingTxsByAddressAndStatusParams, tx *sql.Tx) ([]dbPendingTx.PendingTxsByAddressAndStatusRow, error)
	GetPendingTransactions(ctx context.Context, arg dbPendingTx.GetListPendingTxsParams) ([]dbPendingTx.GetListPendingTxsRow, error)
	GetCountPendingTransaction(ctx context.Context) (int64, error)

	// ---------------- Pending Tx Data ----------------
	InsertPendingTxData(ctx context.Context, args dbPendingTx.InsertPendingTxDataParams, tx *sql.Tx) (dbPendingTx.PendingTxDatum, error)
	FindTxPendingByAddr(ctx context.Context, args dbPendingTx.PendingTxsByAddressParams, tx *sql.Tx) ([]dbPendingTx.PendingTxsByAddressRow, error)
	CountPendingTxByAddr(ctx context.Context, address string, tx *sql.Tx) (int64, error)
	UpdatePendingTxsStatus(ctx context.Context, args dbPendingTx.UpdatePendingTxsStatusParams, tx *sql.Tx) (int64, error)
}

type DbUTXORepository interface {
	GetUTXOByTxIDAndOut(ctx context.Context, arg dbutxo.GetUTXOByTxIDAndOutParams, tx *sql.Tx) (dbutxo.Utxo, error)
	CreateUTXO(ctx context.Context, args dbutxo.CreateUTXOParams, tx *sql.Tx) (dbutxo.Utxo, error)
	DeleteUTXO(ctx context.Context, args dbutxo.DeleteUTXOParams, tx *sql.Tx) error
	DeleteUTXOByBlockID(ctx context.Context, b_id string, tx *sql.Tx) error
	FindUTXOs(ctx context.Context, pubKeyHash string, tx *sql.Tx) ([]dbutxo.Utxo, error)
}
