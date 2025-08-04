package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
)

type Block struct {
	Timestamp    int64          `json:"Timestamp"`
	Hash         []byte         `json:"Hash"`
	PrevHash     []byte         `json:"PrevHash"`
	Transactions []*Transaction `json:"Transactions"`
	Nonce        int64          `json:"Nonce"`
	Height       int64          `json:"Height"`
	MerkleRoot   []byte         `json:"MerkleRoot"`
	Difficulty   int64          `json:"Difficulty"`
	TxCount      int64          `json:"TxCount"`

	NChainWork *big.Int `json:"NChainWork"`
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int64, difficulty int64) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: txs,
		Difficulty:   difficulty,
		Height:       height,
		TxCount:      int64(len(txs)),
	}

	pow := NewProof(block)
	start := time.Now()
	nonce, hash := pow.Run()
	duration := time.Since(start)
	log.Infof("⛏️  Mined block in %s | Nonce: %d | Hash: %x", duration, nonce, hash)

	block.Hash = hash
	block.Nonce = nonce
	block.MerkleRoot = block.HashTransactions()

	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.ErrorHandle(err)
	return res.Bytes()
}

func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.ErrorHandle(err)
	return &block
}

func Genesis(MinerTx *Transaction) *Block {
	return CreateBlock([]*Transaction{MinerTx}, []byte{}, 1, 0)
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serializer())
	}

	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func (b *Block) IsBlockValid(oldBlock Block) bool {
	if b.Height != oldBlock.Height+1 {
		log.Warning("Block is not valid")
		return false
	}

	if b.Timestamp <= oldBlock.Timestamp {
		log.Warning("Invalid timestamp")
		return false
	}

	if !bytes.Equal(oldBlock.Hash, b.PrevHash) {
		log.Warning("Invalid previous hash")
		return false
	}

	if b.Difficulty < MinDifficulty {
		log.Warning("Difficulty too low")
		return false
	}

	if !bytes.Equal(b.MerkleRoot, b.HashTransactions()) {
		log.Warning("Merkle root not match")
		return false
	}

	proof := NewProof(b)
	if !proof.Validate() {
		log.Warning("Proof of Work invalid")
		return false
	}

	return true
}

func (b *Block) IsGenesis() bool {
	return b.PrevHash == nil
}

func ConstructJSON(buffer *bytes.Buffer, b *Block) {
	buffer.WriteString("{")
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%d\",", "Timestamp", b.Timestamp))
	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "PrevHash", b.PrevHash))

	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "Hash", b.Hash))

	buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Difficulty", b.Difficulty))

	buffer.WriteString(fmt.Sprintf("\"%s\":%d,", "Nonce", b.Nonce))

	buffer.WriteString(fmt.Sprintf("\"%s\":\"%x\",", "MerkleRoot", b.MerkleRoot))
	buffer.WriteString(fmt.Sprintf("\"%s\":%d", "TxCount", b.TxCount))
	buffer.WriteString("}")
}

func (b *Block) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	ConstructJSON(buffer, b)
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}
