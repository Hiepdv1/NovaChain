package utils

import (
	"core-blockchain/common/err"
	blockchain "core-blockchain/core"
	"core-blockchain/p2p"
)

type CommandLine struct {
	Blockchain    *blockchain.Blockchain
	P2P           *p2p.Network
	CloseDbAlways bool
}

type BalanceResponse struct {
	Balance   float64
	Address   string
	Timestamp int64
	Error     *err.RPCError
}

type SendResponse struct {
	Message string
	ListTxs []string
	Count   int64
	Error   *err.RPCError
}

type GetMiningTxsResponse struct {
	Message string
	ListTxs []any
	Count   int64
	Error   *err.RPCError
}

type GetBlockResponse struct {
	Block *blockchain.Block
	Error *err.RPCError
}

type GetAllUTXOsResponse struct {
	Message string
	Data    map[string]blockchain.TxOutputs
	Count   int64
	Error   *err.RPCError
}
