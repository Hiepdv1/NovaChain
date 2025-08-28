package wallet

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/env"
	"encoding/json"
)

type walletRPCRepository struct {
	env *env.Env
}

func NewRPCWalletRepository() RPCWalletRepository {
	return &walletRPCRepository{
		env: env.Cfg,
	}
}

func (r *walletRPCRepository) GetBalance(address string) (*Balance, error) {

	params := []any{
		map[string]string{
			"address": address,
		},
	}

	data, err := client.CallRPC(
		r.env.Fullnode_RPC_URL,
		"API.GetBalance",
		params,
	)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	var balance Balance
	if err := json.Unmarshal(rpcResp.Result, &balance); err != nil {
		return nil, err
	}

	return &balance, nil

}
