package chain

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/utils"
)

type Block struct {
	Timestamp    int64              `json:"Timestamp"`
	Hash         string             `json:"Hash"`
	PrevHash     string             `json:"PrevHash"`
	Transactions []*dto.Transaction `json:"Transactions"`
	Nonce        int64              `json:"Nonce"`
	Height       int64              `json:"Height"`
	MerkleRoot   string             `json:"MerkleRoot"`
	NBits        int64              `json:"nBits"`
	TxCount      int64              `json:"TxCount"`
	NChainWork   utils.BigInt       `json:"NChainWork"`
}
