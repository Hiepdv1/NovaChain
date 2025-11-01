package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"core-blockchain/common/err"
	"core-blockchain/json-rpc/types"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type HandleFunc func(parmas json.RawMessage) (any, *err.RPCError)

type API struct {
	cmd *utils.CommandLine
}

var handlers = map[string]HandleFunc{}

func (api *API) HandleCreateWallet(params json.RawMessage) (any, *err.RPCError) {

	wallet, e := api.cmd.CreateWallet()
	if e != nil {
		return nil, err.ErrInternal("Internal Error")
	}

	return wallet, nil
}

func (api *API) HandleGetBalance(params json.RawMessage) (any, *err.RPCError) {
	var args []types.WalletAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	result := api.cmd.GetBalance(args[0].Address)

	return result, nil
}

func (api *API) HandleGetBlockchain(params json.RawMessage) (any, *err.RPCError) {

	var args []types.GetBlockchainAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	blockchain, e := api.cmd.GetBlockChain([]byte(args[0].StartHash), uint16(args[0].Max))
	if e != nil {
		log.Error(e)
		return nil, err.ErrInternal("Internal Error")
	}

	return blockchain, nil
}

func (api *API) HandleGetBlockByHeight(params json.RawMessage) (any, *err.RPCError) {

	var args []types.GetAPIBlockArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	block, e := api.cmd.GetBlockByHeight(args[0].Height)
	if e != nil {
		log.Error(e)
		return nil, err.ErrInternal("Internal error")
	}

	return block, nil
}

func (api *API) HandleGetBlocksByHeightRange(params json.RawMessage) (any, *err.RPCError) {

	var args []types.GetAPIBlockByHeightRangeArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	blocks, e := api.cmd.GetBlocksByHeightRange(args[0].Height, args[0].Limit)

	if e != nil {
		log.Error(e)
		return nil, err.ErrInternal("Internal error")
	}

	return blocks, nil
}

func (api *API) HandleGetCommonBlock(params json.RawMessage) (any, *err.RPCError) {
	var args []types.CommonBlockArgs

	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	block, e := api.cmd.GetCommonBlock(args[0].Locator)
	if e != nil {
		log.Error(e)
		return nil, err.ErrInternal("Internal error")
	}

	if block == nil {
		return nil, err.ErrNotFound("Not common block.")
	}

	return block, nil
}

func (api *API) HandleSendTx(params json.RawMessage) (any, *err.RPCError) {

	var args []types.SendTxAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.SendTx(args[0].TXS), nil
}

func (api *API) GetMiningTxs(params json.RawMessage) (any, *err.RPCError) {

	var args []types.GetMiningTxsAPIArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetMiningTxs(args[0].Verbose), nil
}

func (api *API) GetBlockByHash(params json.RawMessage) (any, *err.RPCError) {

	var args []types.GETBlockByHashArgs
	if e := json.Unmarshal(params, &args); e != nil || len(args) != 1 {
		log.Error(e)
		return nil, err.ErrInvalidArgument("Invalid parameters")
	}

	return api.cmd.GetBlockByHash(args[0].Hash), nil
}

func (api *API) GetAllUTXOs(params json.RawMessage) (any, *err.RPCError) {
	return api.cmd.GetAllUTXOs(), nil
}
