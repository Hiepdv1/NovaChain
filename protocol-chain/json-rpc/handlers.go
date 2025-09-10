package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"core-blockchain/common/err"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type HandleFunc func(parmas json.RawMessage) (any, *err.RPCError)

type API struct {
	cmd *utils.CommandLine
}

var handlers = map[string]HandleFunc{}

func (api *API) HandleCreateWallet(params json.RawMessage) (any, *err.RPCError) {
	log.Info("HandleCreateWallet Executed...")

	return api.cmd.CreateWallet(), nil
}

func (api *API) HandleGetBalance(params json.RawMessage) (any, *err.RPCError) {
	var args []WalletAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	result := api.cmd.GetBalance(args[0].Address)

	return result, nil
}

func (api *API) HandleGetBlockchain(params json.RawMessage) (any, *err.RPCError) {
	log.Info("HandleGetBlockchain Executed...")

	var args []GetBlockchainAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetBlockChain([]byte(args[0].StartHash), uint16(args[0].Max)), nil
}

func (api *API) HandleGetBlockByHeight(params json.RawMessage) (any, *err.RPCError) {
	log.Info("HandleGetBlockByHeight Executed...")

	var args []GetAPIBlockArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetBlockByHeight(args[0].Height), nil
}

func (api *API) HandleGetBlocksByHeightRange(params json.RawMessage) (any, *err.RPCError) {
	log.Info("HandleGetBlocksByHeightRange Executed...")

	var args []GetAPIBlockByHeightRangeArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetBlocksByHeightRange(args[0].Height, args[0].Limit), nil
}

func (api *API) HandleSendTx(params json.RawMessage) (any, *err.RPCError) {
	log.Info("HandleSendTx Executed...")

	var args []SendTxAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return nil, nil

	// return api.cmd.Send(args[0].SendFrom, args[0].SendTo, args[0].Amount, args[0].Fee, args[0].Mine), nil
}

func (api *API) GetBlockByHash(params json.RawMessage) (any, *err.RPCError) {
	log.Info("GetBLockByHas Executed...")

	var args []GETAPIBlockByHash
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetBlockByHash(args[0].Hash), nil
}

func (api *API) GetAllUTXOs(params json.RawMessage) (any, *err.RPCError) {
	log.Info("GetAllUTXOs Executed...")
	return api.cmd.GetAllUTXOs(), nil
}
