package wallet

import (
	"ChainServer/internal/common/client"
)

type rpcWalletRepository struct {
	fullnode_rpc_url string
}

func NewRPCWalletRepository(path string) RPCWalletRepository {
	return &rpcWalletRepository{
		fullnode_rpc_url: path,
	}
}

func (r *rpcWalletRepository) GetBalance(address string) ([]byte, error) {

	params := []interface{}{
		map[string]string{
			"address": address,
		},
	}

	return client.CallRPC(
		r.fullnode_rpc_url,
		"API.GetBalance",
		params,
	)

}
