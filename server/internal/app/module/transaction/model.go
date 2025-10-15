package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/utils"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	dbutxo "ChainServer/internal/db/utxo"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TransactionPending struct {
	dbPendingTx.PendingTransaction
	dbPendingTx.PendingTxDatum
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value      float64
	PubKeyHash []byte
}

type TxInput struct {
	ID        []byte
	Out       int64
	Signature []byte
	PubKey    []byte
}

type TxInputWithDataToSign struct {
	ID         []byte `json:"id"`
	Out        int64  `json:"out"`
	Signature  []byte `json:"signature"`
	PubKey     []byte `json:"pubKey"`
	DataToSign string `json:"dataToSign"`
}

type TransactionWithSigning struct {
	ID      []byte                  `json:"id"`
	Inputs  []TxInputWithDataToSign `json:"inputs"`
	Outputs []TxOutput              `json:"outputs"`
}

func (tx *Transaction) Serializer() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panicf("Transacion - Serializer GobEncode error: %v", err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) TrimmedCopy() Transaction {

	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{
			ID:        in.ID,
			Out:       in.Out,
			Signature: nil,
			PubKey:    nil,
		})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{
			Value:      out.Value,
			PubKeyHash: out.PubKeyHash,
		})
	}

	txCopy := Transaction{
		ID:      tx.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}

	return txCopy
}

func (tx *Transaction) SerializeAndHexEncode() string {

	serializedBytes, err := json.Marshal(tx)
	if err != nil {
		log.Panicf("Failed to JSON marshal transaction: %v", err.Error())
	}

	hexString := fmt.Sprintf("%x", serializedBytes)

	data := utils.DoubleSHA256([]byte(hexString))

	return hex.EncodeToString(data)
}

func (tx *Transaction) BalanceCheck(prevTXs map[string]dbutxo.Utxo) bool {
	totalInput := 0.0
	totalOutput := 0.0

	for _, in := range tx.Inputs {
		prevTx, ok := prevTXs[hex.EncodeToString(in.ID)]
		if !ok {
			return false
		}
		value, err := strconv.ParseFloat(prevTx.Value, 64)
		if err != nil {
			log.Error("Balance Check Error: ", err)
			return false
		}
		totalInput += value
	}

	for _, out := range tx.Outputs {
		totalOutput += out.Value
	}

	fee := totalInput - totalOutput
	if fee < 0 {
		return false
	}

	return totalInput >= totalOutput
}

func (tx *Transaction) WithSigning(prevTXs map[string]dbutxo.Utxo) (*TransactionWithSigning, *apperror.AppError) {

	var inputs []TxInputWithDataToSign

	txCopy := tx.TrimmedCopy()

	balanceOk := tx.BalanceCheck(prevTXs)

	if !balanceOk {
		return nil, apperror.BadRequest("Transaction inputs and outputs do not match", nil)
	}

	for inID, in := range tx.Inputs {
		prevTx, ok := prevTXs[hex.EncodeToString(in.ID)]
		if !ok {
			return nil, apperror.BadRequest("Transaction not found", nil)
		}

		pubKeyHash, err := hex.DecodeString(prevTx.PubKeyHash)
		if err != nil {
			return nil, apperror.Internal("Something went wrong, please try again", nil)
		}

		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = pubKeyHash

		dataToSign := txCopy.SerializeAndHexEncode()

		TxID, err := hex.DecodeString(prevTx.TxID)
		if err != nil {
			return nil, apperror.Internal("Something went wrong, please try again", nil)
		}

		inputs = append(inputs, TxInputWithDataToSign{
			ID:         TxID,
			Out:        prevTx.OutputIndex,
			Signature:  nil,
			PubKey:     in.PubKey,
			DataToSign: dataToSign,
		})

		txCopy.Inputs[inID].PubKey = nil
	}

	txWithSigning := TransactionWithSigning{
		ID:      tx.ID,
		Inputs:  inputs,
		Outputs: tx.Outputs,
	}

	return &txWithSigning, nil
}
