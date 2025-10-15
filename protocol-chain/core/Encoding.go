package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"encoding/binary"
)

func SerializeTransaction(tx *Transaction, buf *bytes.Buffer) {
	utils.WriteBytes(buf, tx.ID)

	binary.Write(buf, binary.LittleEndian, uint32(len(tx.Inputs)))
	for _, in := range tx.Inputs {
		utils.WriteBytes(buf, in.ID)
		binary.Write(buf, binary.LittleEndian, in.Out)
		utils.WriteBytes(buf, in.Signature)
		utils.WriteBytes(buf, in.PubKey)
	}

	binary.Write(buf, binary.LittleEndian, uint32(len(tx.Outputs)))
	for _, out := range tx.Outputs {
		binary.Write(buf, binary.LittleEndian, out.Value)
		utils.WriteBytes(buf, out.PubKeyHash)
	}
}

func SerializeBlock(b *Block) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, b.Timestamp)
	utils.WriteBytes(buf, b.Hash)
	utils.WriteBytes(buf, b.PrevHash)
	utils.WriteBytes(buf, b.MerkleRoot)
	binary.Write(buf, binary.LittleEndian, b.Nonce)
	binary.Write(buf, binary.LittleEndian, b.Height)
	binary.Write(buf, binary.LittleEndian, b.NBits)
	binary.Write(buf, binary.LittleEndian, b.TxCount)
	utils.WriteBigInt(buf, b.NChainWork)

	binary.Write(buf, binary.LittleEndian, uint32(len(b.Transactions)))
	for _, tx := range b.Transactions {
		SerializeTransaction(tx, buf)
	}

	return buf.Bytes()
}

func DeserializeTxData(buf *bytes.Buffer) *Transaction {
	tx := &Transaction{}
	tx.ID = utils.ReadBytes(buf)

	var inCount uint32
	binary.Read(buf, binary.LittleEndian, &inCount)
	for i := uint32(0); i < inCount; i++ {
		in := TxInput{
			ID:        utils.ReadBytes(buf),
			Out:       0,
			Signature: []byte{},
			PubKey:    []byte{},
		}
		binary.Read(buf, binary.LittleEndian, &in.Out)
		in.Signature = utils.ReadBytes(buf)
		in.PubKey = utils.ReadBytes(buf)
		tx.Inputs = append(tx.Inputs, in)
	}

	var outCount uint32
	binary.Read(buf, binary.LittleEndian, &outCount)
	for i := uint32(0); i < outCount; i++ {
		out := TxOutput{}
		binary.Read(buf, binary.LittleEndian, &out.Value)
		out.PubKeyHash = utils.ReadBytes(buf)
		tx.Outputs = append(tx.Outputs, out)
	}

	return tx
}

func DeserializeBlockData(data []byte) *Block {
	buf := bytes.NewBuffer(data)
	b := &Block{}

	binary.Read(buf, binary.LittleEndian, &b.Timestamp)
	b.Hash = utils.ReadBytes(buf)
	b.PrevHash = utils.ReadBytes(buf)
	b.MerkleRoot = utils.ReadBytes(buf)
	binary.Read(buf, binary.LittleEndian, &b.Nonce)
	binary.Read(buf, binary.LittleEndian, &b.Height)
	binary.Read(buf, binary.LittleEndian, &b.NBits)
	binary.Read(buf, binary.LittleEndian, &b.TxCount)
	b.NChainWork = utils.ReadBigInt(buf)

	var txCount uint32
	binary.Read(buf, binary.LittleEndian, &txCount)
	for i := uint32(0); i < txCount; i++ {
		tx := DeserializeTxData(buf)
		b.Transactions = append(b.Transactions, tx)
	}

	return b
}
