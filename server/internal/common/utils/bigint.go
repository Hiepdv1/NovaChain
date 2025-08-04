package utils

import (
	"encoding/json"
	"errors"
	"math/big"
)

type BigInt struct {
	*big.Int
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	bInt := big.NewInt(0)
	_, ok := bInt.SetString(str, 10)
	if !ok {
		return errors.New("invalid big.Int string")
	}

	b.Int = bInt

	return nil
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return json.Marshal("0")
	}

	return json.Marshal(b.String())
}
