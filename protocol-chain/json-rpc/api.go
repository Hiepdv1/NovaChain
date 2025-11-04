package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"core-blockchain/common/err"

	log "github.com/sirupsen/logrus"
)

func NewAPI(cli *utils.CommandLine) *API {
	api := &API{
		cmd: cli,
	}

	api.registerHandlers()

	return api
}

func (api *API) registerHandlers() {
	handlers = map[string]HandleFunc{
		"API.GetAllUTXOs":           api.GetAllUTXOs,
		"API.SendTx":                api.HandleSendTx,
		"API.GetMiningTxs":          api.GetMiningTxs,
		"API.GetBlock":              api.GetBlockByHash,
		"API.GetBalance":            api.HandleGetBalance,
		"API.CreateWallet":          api.HandleCreateWallet,
		"API.GetBlockchain":         api.HandleGetBlockchain,
		"API.GetCommonBlock":        api.HandleGetCommonBlock,
		"API.GetBlockByHeight":      api.HandleGetBlockByHeight,
		"API.GetBlockByHeightRange": api.HandleGetBlocksByHeightRange,
	}
}

func (api *API) ProcessRequest(req JSONRPCRequest) JSONRPCResponse {
	res := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error:   nil,
	}

	handler, ok := handlers[req.Method]

	if !ok {
		res.Error = err.ErrNotFound("Method not found")
	}

	result, e := handler(req.Params)

	if e != nil {
		res.Error = e
	} else {
		b, e := SafeMarshalJSON(result)

		if e != nil {
			log.Error(e)
			res.Error = err.ErrInternal("Encode Error")
			return res
		}

		res.Result = b
	}

	return res
}
