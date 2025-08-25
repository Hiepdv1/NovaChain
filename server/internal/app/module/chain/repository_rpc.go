package chain

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/env"
	"encoding/json"
)

type rpcChainRepository struct {
	env *env.Env
}

func NewRPCChainRepository() RPCChainRepository {
	return &rpcChainRepository{
		env: env.Cfg,
	}
}

func (r *rpcChainRepository) GetBlocks(startHash string, limit int) ([]*Block, error) {

	params := []any{
		map[string]any{
			"startHash": startHash,
			"max":       limit,
		},
	}

	data, err := client.CallRPC(
		r.env.Fullnode_RPC_URL,
		"API.GetBlockchain",
		params,
	)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	var blocks []*Block
	if err := json.Unmarshal(rpcResp.Result, &blocks); err != nil {
		return nil, err
	}

	return blocks, nil

}

func (r *rpcChainRepository) GetBlocksByHeightRange(height, limit int64) ([]*Block, error) {
	params := []any{
		map[string]any{
			"height": height,
			"limit":  limit,
		},
	}

	data, err := client.CallRPC(
		r.env.Fullnode_RPC_URL,
		"API.GetBlockByHeightRange",
		params,
	)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	var blocks []*Block
	if err := json.Unmarshal(rpcResp.Result, &blocks); err != nil {
		return nil, err
	}

	return blocks, nil
}
