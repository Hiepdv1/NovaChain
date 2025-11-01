package blockchain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := CompactToBig(b.NBits)

	pow := &ProofOfWork{b, target}

	log.Infof("Target: %x", target)
	log.Infof("NBits: %d", b.NBits)

	return pow
}

func (pow *ProofOfWork) InitData(nonce int64) ([]byte, error) {
	hashTx, err := pow.Block.HashTransactions()
	if err != nil {
		return nil, err
	}

	info := bytes.Join([][]byte{
		hashTx,
		pow.Block.PrevHash,
		ToByte(int64(nonce)),
		ToByte(int64(pow.Block.NBits)),
		ToByte(pow.Block.Height),
		ToByte(pow.Block.Timestamp),
		ToByte(pow.Block.TxCount),
	}, []byte{})

	return info, nil
}

func (pow *ProofOfWork) Validate() bool {
	var initHash big.Int
	var hash [32]byte

	info, err := pow.InitData(pow.Block.Nonce)
	if err != nil {
		return false
	}
	hash = sha256.Sum256(info)

	initHash.SetBytes(hash[:])

	return initHash.Cmp(pow.Target) == -1

}

func (pow *ProofOfWork) Run(ctx context.Context) (int64, []byte, error) {
	var initHash big.Int
	var hash [32]byte
	var nonce int64

	for nonce = range math.MaxInt64 {
		select {
		case <-ctx.Done():
			return 0, nil, fmt.Errorf("POW: mining stopped manually or context canceled")
		default:
			info, err := pow.InitData(nonce)
			if err != nil {
				return 0, nil, fmt.Errorf("POW: failed to initialize data for nonce %d: %w", nonce, err)
			}

			hash = sha256.Sum256(info)
			initHash.SetBytes(hash[:])

			if initHash.Cmp(pow.Target) == -1 {
				log.Infoln("---------------- Found! ----------------")
				log.Infof("POW: valid hash found! Nonce=%d, Hash=%x", nonce, hash)
				return nonce, hash[:], nil
			}
		}
	}

	return 0, nil, fmt.Errorf("POW: reached max nonce (%d) without finding a valid hash", math.MaxInt64)
}

func ToByte(num int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		log.Error(err)
		return nil
	}

	return buff.Bytes()
}
