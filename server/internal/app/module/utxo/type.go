package utxo

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/types"
)

type TxOutputs struct {
	Outputs []types.TxOutput
}

type GetAllUTXOsRPC struct {
	Message string
	Data    map[string]TxOutputs
	Count   int64
	Error   *client.RPCError
}
