package blockchain

import "math/big"

const (
	PER_COIN             int64 = 100_000_000
	INITIAL_BLOCK_REWARD int64 = 50 * PER_COIN
	HALVING_INTERVAL     int64 = 210000
	MAX_HALVING          int64 = 64
)

func (bc *Blockchain) GetBlockReward(height int64, address string) *Transaction {
	numHalvings := height / HALVING_INTERVAL
	var rewardBlock float64

	if numHalvings >= MAX_HALVING {
		rewardBlock = 0
	} else {
		reward := big.NewInt(INITIAL_BLOCK_REWARD)
		divisor := big.NewInt(1)
		divisor.Lsh(divisor, uint(numHalvings))

		reward.Div(reward, divisor)

		rewardBlock = float64(reward.Int64() / PER_COIN)
	}

	txIn := TxInput{
		[]byte{},
		-1,
		nil,
		[]byte{},
	}

	txOut := NewTxOutput(rewardBlock, address)

	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{*txOut},
	}

	tx.ID = tx.Hash()

	return &tx

}
