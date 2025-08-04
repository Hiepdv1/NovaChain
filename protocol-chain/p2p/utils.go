package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

func GobEncode(data any) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)

	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func BlockForNetwork(block blockchain.Block) blockchain.Block {
	return blockchain.Block{
		Timestamp:    block.Timestamp,
		Height:       block.Height,
		Hash:         block.Hash,
		PrevHash:     block.PrevHash,
		Nonce:        block.Nonce,
		Transactions: block.Transactions,
		MerkleRoot:   block.MerkleRoot,
		Difficulty:   block.Difficulty,
		TxCount:      block.TxCount,
	}
}

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCmd(data []byte) string {
	var cmd []byte
	for _, b := range data {
		if b != byte(0) {
			cmd = append(cmd, b)
		}
	}

	return string(cmd)
}
