package transaction

import (
	"ChainServer/internal/common/client"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"encoding/json"
	"fmt"
)

type rpcTransactionRepository struct {
	env *env.Env
}

func NewRPCTransactionRepo() RpcTransactionRepository {
	return &rpcTransactionRepository{
		env: env.Cfg,
	}
}

func (r *rpcTransactionRepository) SendTx(txs []Transaction) (*RPCSendTxResponse, error) {

	params := []any{
		map[string]any{
			"transactions": txs,
		},
	}

	data, err := client.CallRPC(
		r.env.Fullnode_RPC_URL,
		"API.SendTx",
		params,
	)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("%s", rpcResp.Error.Message)
	}

	var res RPCSendTxResponse
	if err := json.Unmarshal(rpcResp.Result, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func fetchMiningTxs[T any](url string, params []any) (*RPCGetMiningTxResponse[T], error) {
	data, err := client.CallRPC(url, "API.GetMiningTxs", params)
	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, err
	}

	var res RPCGetMiningTxResponse[T]
	if err := json.Unmarshal(rpcResp.Result, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *rpcTransactionRepository) FetchMiningTxIDs() (*RPCGetMiningTxResponse[string], error) {
	params := []any{map[string]any{"verbose": false}}
	return fetchMiningTxs[string](r.env.Fullnode_RPC_URL, params)
}

func (r *rpcTransactionRepository) FetchMiningTxsFull() (*RPCGetMiningTxResponse[dto.Transaction], error) {
	params := []any{map[string]any{"verbose": true}}
	return fetchMiningTxs[dto.Transaction](r.env.Fullnode_RPC_URL, params)
}
