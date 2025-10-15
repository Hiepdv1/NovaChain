package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

func WriteBytes(buf *bytes.Buffer, data []byte) {
	if data == nil {
		binary.Write(buf, binary.LittleEndian, uint32(0))
		return
	}
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
}

func WriteBigInt(buf *bytes.Buffer, n *big.Int) {
	if n == nil {
		WriteBytes(buf, nil)
		return
	}
	WriteBytes(buf, n.Bytes())
}

func ReadBytes(buf *bytes.Buffer) []byte {
	var length uint32
	binary.Read(buf, binary.LittleEndian, &length)
	if length == 0 {
		return nil
	}
	data := make([]byte, length)
	buf.Read(data)
	return data
}

func ReadBigInt(buf *bytes.Buffer) *big.Int {
	data := ReadBytes(buf)
	if data == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(data)
}
