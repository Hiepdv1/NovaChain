package utils

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/types"
	dbutxo "ChainServer/internal/db/utxo"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func newTxOutput(amount float64, to []byte) types.TxOutputBytes {
	pubKeyHash := to[1 : len(to)-checkSumlength]

	output := types.TxOutputBytes{Value: amount, PubKeyHash: pubKeyHash}

	return output
}

func transactionHash(tx *types.TransactionBytes) []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = nil

	hash = sha256.Sum256(txCopy.Serializer())

	return hash[:]
}

func NewTransaction(pubkey []byte, from, to []byte, amount, fee float64, utxos []dbutxo.Utxo, accUtxo float64) (*types.TransactionBytes, *apperror.AppError) {
	PER_COIN := 100_000_000

	if fee < float64(1/PER_COIN) {
		return nil, apperror.BadRequest(fmt.Sprintf("fee must be greater than or equal 1/%d (%f)", PER_COIN, float64(1/PER_COIN)), nil)
	}

	var inputs []types.TxInputBytes
	var outputs []types.TxOutputBytes

	if accUtxo < amount+fee {
		return nil, apperror.BadRequest("you dont have enough amount", nil)
	}

	var currentAcc float64

	for _, utxo := range utxos {
		value, err := strconv.ParseFloat(utxo.Value, 64)
		if err != nil {
			log.Error("New Transaction Error: ", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", nil)
		}
		currentAcc += value

		txID, err := hex.DecodeString(utxo.TxID.String)
		if err != nil {
			log.Error("New Transaction Error: ", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", nil)
		}

		newInput := types.TxInputBytes{
			ID:     txID,
			Out:    utxo.OutputIndex,
			PubKey: pubkey,
		}

		inputs = append(inputs, newInput)
	}

	outputs = append(outputs, newTxOutput(amount, to))

	if currentAcc >= amount+fee {
		outputs = append(outputs, newTxOutput(accUtxo-amount-fee, from))
	}

	tx := types.TransactionBytes{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.ID = transactionHash(&tx)

	return &tx, nil
}
