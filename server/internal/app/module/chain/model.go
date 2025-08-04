package chain

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/utils"
)

type Block struct {
	Timestamp    int64                      `json:"Timestamp"`
	Hash         string                     `json:"Hash"`
	PrevHash     string                     `json:"PrevHash"`
	Transactions []*transaction.Transaction `json:"Transactions"`
	Nonce        int64                      `json:"Nonce"`
	Height       int64                      `json:"Height"`
	MerkleRoot   string                     `json:"MerkleRoot"`
	Difficulty   int64                      `json:"Difficulty"`
	TxCount      int64                      `json:"TxCount"`
	NChainWork   utils.BigInt               `json:"NChainWork"`
}
