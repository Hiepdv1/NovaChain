package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type HandleFunc func(parmas json.RawMessage) (any, *RPCError)

type API struct {
	cmd *utils.CommandLine
}

var handlers = map[string]HandleFunc{}

func (api *API) HandleCreateWallet(params json.RawMessage) (any, *RPCError) {
	log.Info("HandleCreateWallet Executed...")

	return api.cmd.CreateWallet(), nil
}

func (api *API) HandleGetBalance(params json.RawMessage) (any, *RPCError) {
	var args []WalletAPIArgs
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		return nil, &RPCError{
			Code:    -32602,
			Message: "Invalid parameters",
		}
	}

	result := api.cmd.GetBalance(args[0].Address)

	return result, nil
}

func (api *API) HandleGetBlockchain(params json.RawMessage) (any, *RPCError) {
	log.Info("HandleGetBlockchain Executed...")

	var args []GetBlockchainAPIArgs
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		log.Error(err)
		return nil, &RPCError{
			Code:    -32602,
			Message: "Invalid parameters",
		}
	}

	return api.cmd.GetBlockChain([]byte(args[0].StartHash), uint16(args[0].Max)), nil
}

func (api *API) HandleGetBlockByHeight(params json.RawMessage) (any, *RPCError) {
	log.Info("HandleGetBlockByHeight Executed...")

	var args []GetAPIBlockArgs
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		log.Error(err)
		return nil, &RPCError{
			Code:    -32602,
			Message: "Invalid parameters",
		}
	}

	return api.cmd.GetBlockByHeight(args[0].Height), nil
}

func (api *API) HandleGetBlocksByHeightRange(params json.RawMessage) (any, *RPCError) {
	log.Info("HandleGetBlocksByHeightRange Executed...")

	var args []GetAPIBlockByHeightRangeArgs
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		log.Error(err)
		return nil, &RPCError{
			Code:    -32602,
			Message: "Invalid parameters",
		}
	}

	return api.cmd.GetBlocksByHeightRange(args[0].Height, args[0].Limit), nil
}

func (api *API) HandleSendTx(params json.RawMessage) (any, *RPCError) {
	log.Info("HandleSendTx Executed...")

	var args []SendTxAPIArgs
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		log.Error(err)
		return nil, &RPCError{
			Code:    -32602,
			Message: "Invalid parameters",
		}
	}

	return api.cmd.Send(args[0].SendFrom, args[0].SendTo, args[0].Amount, args[0].Fee, args[0].Mine), nil
}
