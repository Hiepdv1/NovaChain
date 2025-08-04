package jsonrpc

import (
	"core-blockchain/cmd/utils"

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
		"API.SendTx":                api.HandleSendTx,
		"API.GetBalance":            api.HandleGetBalance,
		"API.CreateWallet":          api.HandleCreateWallet,
		"API.GetBlockchain":         api.HandleGetBlockchain,
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
		res.Error = &RPCError{
			Code:    -32602,
			Message: "Method not found",
		}
	}

	result, err := handler(req.Params)

	if err != nil {
		res.Error = err
	} else {
		b, err := SafeMarshalJSON(result)

		if err != nil {
			log.Error(err)
			res.Error = &RPCError{
				Code:    -32602,
				Message: "Encode Error",
			}
			return res
		}

		res.Result = b
	}

	return res
}
