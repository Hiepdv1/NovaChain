package types

import (
	"bytes"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	ID      string
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID        string
	Out       int64
	Signature string
	PubKey    string
}

type TxOutput struct {
	Value      float64
	PubKeyHash string
}

type TxInputBytes struct {
	ID        []byte `json:"id"`
	Out       int64  `json:"out"`
	Signature []byte `json:"signature"`
	PubKey    []byte `json:"pubKey"`
}

type TxOutputBytes struct {
	Value      float64 `json:"value"`
	PubKeyHash []byte  `json:"pubKeyHash"`
}

type TransactionBytes struct {
	ID      []byte          `json:"id"`
	Inputs  []TxInputBytes  `json:"inputs"`
	Outputs []TxOutputBytes `json:"outputs"`
}

func (tx *TransactionBytes) Serializer() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panicf("Transacion - Serializer GobEncode error: %v", err)
	}

	return encoded.Bytes()
}
