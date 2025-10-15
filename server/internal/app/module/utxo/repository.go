package utxo

import (
	dbutxo "ChainServer/internal/db/utxo"
	"context"
	"database/sql"
)

type RPCUtxoRepository interface {
	GetAllUTXOs() (*GetAllUTXOsRPC, error)
}

type DbUTXORepository interface {
	GetUTXOByTxIDAndOut(ctx context.Context, arg dbutxo.GetUTXOByTxIDAndOutParams, tx *sql.Tx) (dbutxo.Utxo, error)
	CreateUTXO(ctx context.Context, args dbutxo.CreateUTXOParams, tx *sql.Tx) (dbutxo.Utxo, error)
	DeleteUTXO(ctx context.Context, args dbutxo.DeleteUTXOParams, tx *sql.Tx) error
	DeleteUTXOByBlockID(ctx context.Context, b_id string, tx *sql.Tx) error
	FindUTXOs(ctx context.Context, pubKeyHash string, tx *sql.Tx) ([]dbutxo.Utxo, error)
}
