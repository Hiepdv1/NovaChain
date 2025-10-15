package blockchain

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	NBits        uint32         `json:"nBits"`
	TxCount      int64          `json:"TxCount"`

	NChainWork *big.Int `json:"NChainWork"`
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int64, NBits uint32, ctx context.Context) (*Block, error) {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: txs,
		NBits:        NBits,
		Height:       height,
		TxCount:      int64(len(txs)),
	}

	pow := NewProof(block)
	start := time.Now()
	nonce, hash, err := pow.Run(ctx)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	log.Infof("⛏️  Mined block in %s | Nonce: %d | Hash: %x", duration, nonce, hash)

	block.Hash = hash
	block.Nonce = *nonce

	merkleRoot, err := block.HashTransactions()
	if err != nil {
		return nil, err
	}

	block.MerkleRoot = merkleRoot

	return block, nil
}

func (b *Block) Size() (int, error) {

	dataBytes := SerializeBlock(b)

	return len(dataBytes), nil
}

func Genesis(MinerTx *Transaction) (*Block, error) {
	block := &Block{
		Timestamp:    1758441999,
		PrevHash:     nil,
		Transactions: []*Transaction{MinerTx},
		NBits:        0x207fffff,
		Height:       1,
		TxCount:      1,
	}

	pow := NewProof(block)
	start := time.Now()
	nonce := int64(1)
	info, err := pow.InitData(nonce)

	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(info)

	duration := time.Since(start)
	log.Infof("---> Mined block in %s | Nonce: %d | Hash: %x <---", duration, nonce, hash)

	block.Hash = hash[:]
	block.Nonce = nonce

	merkleRoot, err := block.HashTransactions()
	if err != nil {
		return nil, err
	}

	block.MerkleRoot = merkleRoot

	return block, nil
}

func (b *Block) HashTransactions() ([]byte, error) {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		serialize := new(bytes.Buffer)
		SerializeTransaction(tx, serialize)

		txHashes = append(txHashes, serialize.Bytes())
	}

	tree, err := NewMerkleTree(txHashes)
	if err != nil {
		return nil, err
	}

	return tree.RootNode.Data, nil
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

	merkleRoot, err := b.HashTransactions()
	if err != nil {
		log.Error(err)
		return false
	}

	if !bytes.Equal(b.MerkleRoot, merkleRoot) {
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
