package blockchain

import (
	"bytes"
	"core-blockchain/common/env"
	"core-blockchain/wallet"
	"encoding/gob"
)

var conf = env.New()

var (
	checkSumlength = conf.WalletAddressCheckSum
	// version        = byte(0x00)
)

type TxInput struct {
	ID        []byte
	Out       int64
	Signature []byte
	PubKey    []byte
}

type TxOutput struct {
	Value      float64
	PubKeyHash []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

func NewTxOutput(value float64, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : int64(len(pubKeyHash))-checkSumlength]

	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

func (outs *TxOutputs) Serialize() ([]byte, error) {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(outs)

	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil

}

func DeSerializeOuputs(data []byte) (*TxOutputs, error) {
	var outputs TxOutputs

	encoder := gob.NewDecoder(bytes.NewReader(data))

	err := encoder.Decode(&outputs)

	if err != nil {
		return nil, err
	}

	return &outputs, nil
}
