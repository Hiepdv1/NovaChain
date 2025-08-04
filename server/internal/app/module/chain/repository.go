package chain

import (
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
)

type RPCChainRepository interface {
	GetBlocks(startHash string, limit int) ([]*Block, error)
	GetBlocksByHeightRange(height, limit int64) ([]*Block, error)
}

type DBChainRepository interface {
	CreateBlock(ctx context.Context, args dbchain.CreateBlockParams, tx *sql.Tx) (dbchain.Block, error)
	GetBlockByHeight(ctx context.Context, height int64) (dbchain.Block, error)
	GetLastBlock(ctx context.Context) (dbchain.Block, error)
	GetBlockByHash(ctx context.Context, hash string) (dbchain.Block, error)
	DeleteBlockByHash(ctx context.Context, hash string, tx *sql.Tx) error
	GetListBlock(ctx context.Context, args dbchain.GetListBlocksParams) ([]dbchain.Block, error)
}
