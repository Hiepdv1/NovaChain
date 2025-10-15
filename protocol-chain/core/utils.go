package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

func GobEncode[T any](value T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return nil, fmt.Errorf("failed to encode value: %v", err)
	}

	return buf.Bytes(), nil
}

func GobDecode[T any](data []byte) (T, error) {
	var result T

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode value: %v", err)
	}

	return result, nil
}

func DoubleSHA256(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:]
}
