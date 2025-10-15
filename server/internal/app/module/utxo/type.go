package utxo

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/types"
)

type GetAllUTXOsRPC struct {
	Message string
	Data    map[string]types.TxOutputs
	Count   int64
	Error   *client.RPCError
}
