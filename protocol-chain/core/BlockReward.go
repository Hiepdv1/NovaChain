package blockchain

import "math/big"

const (
	PER_COIN             int64 = 100_000_000
	INITIAL_BLOCK_REWARD int64 = 50 * PER_COIN
	HALVING_INTERVAL     int64 = 210000
	MAX_HALVING          int64 = 64
)

func (bc *Blockchain) GetReward(height int64) *CoinAmount {
	numHalvings := height / HALVING_INTERVAL

	if numHalvings >= MAX_HALVING {
		return ZeroAmount()
	}

	reward := big.NewInt(INITIAL_BLOCK_REWARD)

	divisor := big.NewInt(1)
	divisor.Lsh(divisor, uint(numHalvings))

	reward.Div(reward, divisor)

	return &CoinAmount{value: reward}
}

func (bc *Blockchain) GetBlockReward(height int64, address string) (*Transaction, error) {
	rewardBlock := bc.GetReward(height)

	txIn := TxInput{
		[]byte{},
		-1,
		[]byte{},
		[]byte{},
	}

	txOut := NewTxOutput(rewardBlock.ToFloat(), address)

	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{*txOut},
	}

	txIdHash, err := tx.Hash(height)
	if err != nil {
		return nil, err
	}

	tx.ID = txIdHash

	return &tx, nil

}
