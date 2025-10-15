package blockchain

import (
	"bytes"
	"core-blockchain/wallet"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func NewTransaction(w *wallet.Wallet, to string, amount, fee float64, utxo *UTXOSet, height int64) (*Transaction, error) {
	if fee < float64(1/PER_COIN) {
		return nil, fmt.Errorf("fee must be greater than or equal 1/%d (%f)", PER_COIN, float64(1/PER_COIN))
	}

	var inputs []TxInput
	var outputs []TxOutput

	publicKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs, err := utxo.FindSpendableOutputs(publicKeyHash, amount+fee)
	if err != nil {
		return nil, err
	}
	if acc < amount+fee {
		err := errors.New("you dont have enough amount")
		return nil, err
	}

	from := string(w.Address())

	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)

		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := TxInput{txID, int64(out), nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTxOutput(amount, to))

	if acc >= amount+fee {
		outputs = append(outputs, *NewTxOutput(acc-amount-fee, from))
	}

	tx := Transaction{nil, inputs, outputs}
	txIdhash, err := tx.Hash(height)

	if err != nil {
		return nil, err
	}

	tx.ID = txIdhash

	utxo.Blockchain.SignTransaction(w.PrivateKey, &tx)

	return &tx, nil
}

func (tx *Transaction) Hash(height int64) ([]byte, error) {
	var hash [32]byte
	buf := new(bytes.Buffer)

	txCopy := *tx
	txCopy.ID = nil

	SerializeTransaction(&txCopy, buf)

	binary.Write(buf, binary.LittleEndian, height)

	hash = sha256.Sum256(buf.Bytes())

	return hash[:], nil
}

func (tx *Transaction) IsMinerTx() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) error {
	if tx.IsMinerTx() {
		return nil
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Fatal("ERROR: Previous Transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]

		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash

		dataBytes, err := json.Marshal(txCopy)
		if err != nil {
			return err
		}

		dataToSign := fmt.Sprintf("%x", dataBytes)

		digest := DoubleSHA256([]byte(dataToSign))

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, digest)

		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
		txCopy.Inputs[inId].PubKey = nil
	}

	return nil
}

func (tx *Transaction) BalanceCheck(prevTXs map[string]Transaction) bool {
	var totalInput float64 = 0
	var totalOutput float64 = 0

	for _, in := range tx.Inputs {
		txID := hex.EncodeToString(in.ID)
		prevTx, ok := prevTXs[txID]
		if !ok {
			return false
		}

		if in.Out < 0 || in.Out >= int64(len(prevTx.Outputs)) {
			return false
		}

		pubKeyHash := wallet.PublicKeyHash(in.PubKey)
		if !bytes.Equal(pubKeyHash, prevTx.Outputs[in.Out].PubKeyHash) {
			return false
		}

		totalInput += prevTx.Outputs[in.Out].Value
	}

	for _, out := range tx.Outputs {
		totalOutput += out.Value
	}

	return totalInput >= totalOutput
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Error("ERROR: Previous Transaction is not valid")
			return false
		}
	}

	if !tx.BalanceCheck(prevTXs) {
		return false
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]

		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash

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

		dataByte, err := json.Marshal(txCopy)
		if err != nil {
			log.Errorf("Failed to JSON marshal transaction: %v", err)
			return false
		}
		dataToVerify := fmt.Sprintf("%x", dataByte)

		digest := DoubleSHA256([]byte(dataToVerify))

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		if !ecdsa.Verify(&rawPubKey, digest, &r, &s) {
			return false
		}

		txCopy.Inputs[inId].PubKey = nil
	}

	return true

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

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("---Transaction: %x", tx.ID))

	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf(" Input (%d): ", i))
		lines = append(lines, fmt.Sprintf(" 	 	TXID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("		Out: %d", input.Out))
		lines = append(lines, fmt.Sprintf(" 	 	Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf(" Output: (%d): ", i))
		lines = append(lines, fmt.Sprintf(" 	 	Value: %f", output.Value))
		lines = append(lines, fmt.Sprintf("		PubkeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

func InitGenesisTx(height int64) (*Transaction, error) {

	pubkey, err := hex.DecodeString("5d5807642aea55229a534a596b0b98c76346abccf85c83d17e7e80cfc9eef4682c85ff581cc4ca8d26244e9d3d9ed0695241c76ac288bb3a9b3b7802da0db4b7")

	if err != nil {
		log.Panicf("Failed to decode pubKey: %v", err)
	}

	txIn := TxInput{
		ID:        []byte{},
		Out:       -1,
		Signature: []byte{},
		PubKey:    pubkey,
	}
	txOut := NewTxOutput(111_111_111.965185, "1LacjauKAjDJA34hjS9xJ2uEez7pQYqh5N")

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
