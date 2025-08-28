package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"core-blockchain/wallet"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

func NewTransaction(w *wallet.Wallet, to string, amount, fee float64, utxo *UTXOSet) (*Transaction, error) {
	if fee < float64(1/PER_COIN) {
		return nil, fmt.Errorf("fee must be greater than or equal 1/%d (%f)", PER_COIN, float64(1/PER_COIN))
	}

	var inputs []TxInput
	var outputs []TxOutput

	publicKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := utxo.FindSpendableOutputs(publicKeyHash, amount+fee)

	if acc < amount+fee {
		err := errors.New("you dont have enough amount")
		return nil, err
	}

	from := string(w.Address())

	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)

		utils.ErrorHandle(err)

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
	tx.ID = tx.Hash()

	utxo.Blockchain.SignTransaction(w.PrivateKey, &tx)

	return &tx, nil
}

func (tx *Transaction) Serializer() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)

	err := encoder.Encode(tx)

	utils.ErrorHandle(err)

	return encoded.Bytes()
}

func (tx *Transaction) Deserialize(data []byte) *Transaction {
	var transaction Transaction

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&transaction)
	utils.ErrorHandle(err)

	return &transaction

}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serializer())

	return hash[:]
}

func (tx *Transaction) IsMinerTx() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsMinerTx() {
		return
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

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		utils.ErrorHandle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
		txCopy.Inputs[inId].PubKey = nil
	}

}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Fatal("ERROR: Previous Transaction is not valid")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]

		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[inId].PubKeyHash

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

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		if !ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) {
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

func InitGenesisTx() *Transaction {

	pubkey, err := hex.DecodeString("35c65154f5407780a2b5e5303fb52d46092a1c4f08ffa8a20c08ed007869af3c9873d5c743aeb610de7ba56fa338ac22a8d863d0f7b67cad3eee653ea2f93e61")

	if err != nil {
		log.Panicf("Failed to decode pubKey: %v", err)
	}

	txIn := TxInput{
		ID:        []byte{},
		Out:       -1,
		Signature: nil,
		PubKey:    pubkey,
	}
	txOut := NewTxOutput(1_111_111, "13eQHPE45d95ayyVoxduzPDt13GWF5is1Q")

	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{*txOut},
	}

	tx.ID = tx.Hash()

	return &tx
}
