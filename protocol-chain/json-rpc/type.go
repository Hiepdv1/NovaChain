package jsonrpc

import (
	"core-blockchain/common/err"
	blockchain "core-blockchain/core"
	"encoding/json"
)

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
	TXS []*blockchain.Transaction `json:"transactions"`
}

type GetMiningTxsAPIArgs struct {
	Verbose bool `json:"verbose"`
}

type GETAPIBlockByHash struct {
	Hash []byte `json:"hash"`
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
	Error   *err.RPCError   `json:"error,omitempty"`
	ID      any             `json:"id"`
}
