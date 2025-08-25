package jsonrpc

import "encoding/json"

type RPCError struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type WalletAPIArgs struct {
	Address string `json:"address"`
}

type GetBlockchainAPIArgs struct {
	StartHash string `json:"startHash"`
	Max       uint64 `json:"max"`
}

type GetAPIBlockArgs struct {
	Height int64 `json:"height"`
}

type GetAPIBlockByHeightRangeArgs struct {
	Height int64 `json:"height"`
	Limit  int64 `json:"limit"`
}

type SendTxAPIArgs struct {
	SendFrom string  `json:"sendFrom"`
	SendTo   string  `json:"sendTo"`
	Amount   float64 `json:"amount"`
	Fee      float64 `json:"fee"`
	Mine     bool    `json:"mine"`
}

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      any             `json:"id"`
}
