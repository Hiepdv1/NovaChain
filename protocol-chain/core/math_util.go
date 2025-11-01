package blockchain

import (
	"fmt"
	"math/big"
	"strings"
)

const CoinPrecision = PER_COIN

type CoinAmount struct {
	value *big.Int
}

func NewCoinAmountFromInt(v int64) *CoinAmount {
	return &CoinAmount{value: big.NewInt(v * int64(CoinPrecision))}
}

func NewCoinAmountFromFloat(v float64) *CoinAmount {
	f := new(big.Float).Mul(big.NewFloat(v), big.NewFloat(float64(CoinPrecision)))
	i := new(big.Int)
	f.Int(i)
	return &CoinAmount{value: i}
}

func NewCoinAmountFromString(s string) (*CoinAmount, error) {
	f, ok := new(big.Float).SetString(s)
	if !ok {
		return nil, fmt.Errorf("invalid numeric string: %s", s)
	}
	f.Mul(f, big.NewFloat(float64(CoinPrecision)))
	i := new(big.Int)
	f.Int(i)
	return &CoinAmount{value: i}, nil
}

func ZeroAmount() *CoinAmount {
	return &CoinAmount{value: big.NewInt(0)}
}

func (c *CoinAmount) Add(b *CoinAmount) *CoinAmount {
	return &CoinAmount{value: new(big.Int).Add(c.value, b.value)}
}

func (c *CoinAmount) Sub(b *CoinAmount) *CoinAmount {
	return &CoinAmount{value: new(big.Int).Sub(c.value, b.value)}
}

func (c *CoinAmount) Mul(factor int64) *CoinAmount {
	return &CoinAmount{value: new(big.Int).Mul(c.value, big.NewInt(factor))}
}

func (c *CoinAmount) Div(divisor int64) *CoinAmount {
	if divisor == 0 {
		return ZeroAmount()
	}
	return &CoinAmount{value: new(big.Int).Div(c.value, big.NewInt(divisor))}
}

func (c *CoinAmount) Cmp(b *CoinAmount) int {
	return c.value.Cmp(b.value)
}

func (c *CoinAmount) IsNegative() bool {
	return c.value.Sign() < 0
}

func (c *CoinAmount) Clone() *CoinAmount {
	return &CoinAmount{value: new(big.Int).Set(c.value)}
}

func (c *CoinAmount) ToFloat() float64 {
	f := new(big.Float).Quo(
		new(big.Float).SetInt(c.value),
		big.NewFloat(float64(CoinPrecision)),
	)
	v, _ := f.Float64()
	return v
}

func (c *CoinAmount) String() string {
	f := new(big.Float).Quo(new(big.Float).SetInt(c.value), big.NewFloat(float64(CoinPrecision)))
	str := f.Text('f', 8)
	return strings.TrimRight(strings.TrimRight(str, "0"), ".")
}

func (c *CoinAmount) Raw() *big.Int {
	return new(big.Int).Set(c.value)
}

func SumFees(totalInput, totalOutput *CoinAmount) *CoinAmount {
	return totalInput.Sub(totalOutput)
}

func ValidateBlockReward(reward, expected *CoinAmount) bool {
	return reward.Cmp(expected) <= 0
}
