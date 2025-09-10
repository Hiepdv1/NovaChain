package utxo

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/env"
	"encoding/json"
)

type rpcUtxoRepository struct {
	env *env.Env
}

func NewRPCUtxoRepository() RPCUtxoRepository {
	return &rpcUtxoRepository{
		env: env.Cfg,
	}
}

func (r *rpcUtxoRepository) GetAllUTXOs() (*GetAllUTXOsRPC, error) {
	params := []any{
		map[string]any{},
	}

	data, err := client.CallRPC(
		r.env.Fullnode_RPC_URL,
		"API.GetAllUTXOs",
		params,
	)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	var utxos GetAllUTXOsRPC
	if err := json.Unmarshal(rpcResp.Result, &utxos); err != nil {
		return nil, err
	}

	return &utxos, nil
}
