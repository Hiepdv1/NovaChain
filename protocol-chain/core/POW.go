package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Difficulty))

	pow := &ProofOfWork{b, target}

	log.Infof("Target: %x\n", target)
	log.Infof("Difficulty: %x\n", b.Difficulty)

	return pow
}

func (pow *ProofOfWork) InitData(nonce int64) []byte {
	info := bytes.Join([][]byte{
		pow.Block.HashTransactions(),
		pow.Block.PrevHash,
		ToByte(int64(nonce)),
		ToByte(int64(pow.Block.Difficulty)),
		ToByte(pow.Block.Height),
		ToByte(pow.Block.Timestamp),
		ToByte(pow.Block.TxCount),
	}, []byte{})

	return info
}

func (pow *ProofOfWork) Validate() bool {
	var initHash big.Int
	var hash [32]byte

	info := pow.InitData(pow.Block.Nonce)
	hash = sha256.Sum256(info)

	initHash.SetBytes(hash[:])

	return initHash.Cmp(pow.Target) == -1

}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var initHash big.Int
	var hash [32]byte
	var nonce int64

	for nonce = range math.MaxInt64 {
		info := pow.InitData(nonce)
		hash = sha256.Sum256(info)

		log.Infof("Pow: \r%x", hash)
		initHash.SetBytes(hash[:])

		if initHash.Cmp(pow.Target) == -1 {
			log.Infoln("---------------- Found! ----------------")
			break
		}
	}

	return nonce, hash[:]
}

func ToByte(num int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)

	utils.ErrorHandle(err)

	return buff.Bytes()
}
