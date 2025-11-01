package types

import blockchain "core-blockchain/core"

type CommonBlockArgs struct {
	Locator [][]byte `json:"locator"`
}

type GETBlockByHashArgs struct {
	Hash []byte `json:"hash"`
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
	TXS []*blockchain.Transaction `json:"transactions"`
}

type GetMiningTxsAPIArgs struct {
	Verbose bool `json:"verbose"`
}
