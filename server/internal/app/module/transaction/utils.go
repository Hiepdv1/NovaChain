package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/env"
	dbutxo "ChainServer/internal/db/utxo"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func newTxOutput(amount float64, to []byte) TxOutput {
	pubKeyHash := to[1 : len(to)-int(env.Cfg.CheckSumLength)]

	output := TxOutput{Value: amount, PubKeyHash: pubKeyHash}

	return output
}

func transactionHash(tx *Transaction) []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = nil

	hash = sha256.Sum256(txCopy.Serializer())

	return hash[:]
}

func NewTransaction(pubkey []byte, from, to []byte, amount, fee float64, utxos []dbutxo.Utxo, accUtxo float64) (*Transaction, *apperror.AppError) {
	PER_COIN := 100_000_000

	if fee < float64(1/PER_COIN) {
		return nil, apperror.BadRequest(fmt.Sprintf("fee must be greater than or equal 1/%d (%f)", PER_COIN, float64(1/PER_COIN)), nil)
	}

	var inputs []TxInput
	var outputs []TxOutput

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

		txID, err := hex.DecodeString(utxo.TxID)
		if err != nil {
			log.Error("New Transaction Error: ", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", nil)
		}

		newInput := TxInput{
			ID:        txID,
			Out:       utxo.OutputIndex,
			Signature: nil,
			PubKey:    pubkey,
		}

		inputs = append(inputs, newInput)

		if currentAcc >= amount+fee {
			break
		}
	}

	outputs = append(outputs, newTxOutput(amount, to))

	if currentAcc >= amount+fee {
		outputs = append(outputs, newTxOutput(accUtxo-amount-fee, from))
	}

	tx := Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.ID = transactionHash(&tx)

	return &tx, nil
}

func VerifyTransactionSig(tx *Transaction, utxos map[string]dbutxo.Utxo) bool {
	curve := elliptic.P256()

	for _, in := range tx.Inputs {
		if _, ok := utxos[hex.EncodeToString(in.ID)]; !ok {
			return false
		}
	}

	txWithSigning, apperr := tx.WithSigning(utxos)

	if apperr != nil {
		return false
	}

	for inID, in := range tx.Inputs {
		r := big.Int{}
		s := big.Int{}
		SigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(SigLen / 2)])
		s.SetBytes(in.Signature[(SigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		dataToVerify := txWithSigning.Inputs[inID].DataToSign
		dataToVerifyBytes, err := hex.DecodeString(dataToVerify)
		if err != nil {
			log.Error(err)
			return false
		}

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		SigOk := ecdsa.Verify(&rawPubKey, dataToVerifyBytes, &r, &s)

		if !SigOk {
			return false
		}
	}

	return true

}
