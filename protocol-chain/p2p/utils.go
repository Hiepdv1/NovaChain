package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
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
		NBits:        block.NBits,
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

func SliceMap[T any, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func LoadOrCreateIdentity(path string) (libp2pcrypto.PrivKey, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("failed to create key directory: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}

		priv, err := libp2pcrypto.UnmarshalPrivateKey(data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal private key: %w", err)
		}

		log.Infof("🔑 Loaded existing identity from: %s", path)
		return priv, nil
	}

	priv, _, err := libp2pcrypto.GenerateKeyPair(libp2pcrypto.Ed25519, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	data, err := libp2pcrypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write private key: %w", err)
	}

	fmt.Println("✅ Created new identity and saved to:", path)
	return priv, nil
}
