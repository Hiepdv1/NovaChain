package transaction

import (
	"ChainServer/internal/common/client"
	dbchain "ChainServer/internal/db/chain"
	"time"
)

type TransactionDataParsed struct {
	Fee       float64
	Amount    float64
	To        []byte // base58
	Timestamp time.Time
	Message   string
}

type NewTransactionParsed struct {
	Data TransactionDataParsed
	Sig  []byte
}

type SendTransactionDataParsed struct {
	Fee          float64
	Amount       float64
	ReceiverAddr string
	Priority     uint
	Transaction  Transaction
}

type RPCSendTxResponse struct {
	Message string
	ListTxs []string
	Count   int64
	Error   *client.RPCError
}

type RPCGetMiningTxResponse[T any] struct {
	Message string
	Txs     []T
	Count   int64
	Error   *client.RPCError
}

type DetailTransaction struct {
	dbchain.GetDetailTxRow
	Difficulty int64
}
