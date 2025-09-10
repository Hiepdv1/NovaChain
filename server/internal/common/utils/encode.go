package utils

import (
	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	enCode := base58.Encode(input)

	return []byte(enCode)
}

func Base58Decode(input string) ([]byte, error) {
	decode, err := base58.Decode(input)

	if err != nil {
		return nil, err
	}

	return decode, nil
}
