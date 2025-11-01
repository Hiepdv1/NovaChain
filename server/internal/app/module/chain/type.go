package chain

import (
	"ChainServer/internal/common/response"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type BlockDetail struct {
	ID           uuid.UUID
	BID          string
	PrevHash     sql.NullString
	Nonce        int64
	Height       int64
	MerkleRoot   string
	Nbits        int64
	TxCount      int64
	NchainWork   string
	Size         float64
	Timestamp    int64
	TotalFee     interface{}
	Difficulty   int64
	Miner        string
	Transactions struct {
		Data json.RawMessage
		Meta *response.PaginationMeta
	}
}

type NetworkInfo struct {
	LastBlock     int64
	Hashrate      string
	AvgBlockTime  float64
	AvgDifficulty float64
	SyncStatus    string
	NetworkHealth string
	TxPending     int64
}
