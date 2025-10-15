package utils

import "math/big"

func CompactToTarget(nBits uint32) *big.Int {
	exponent := byte(nBits >> 24)
	mantissa := nBits & 0x007fffff

	target := new(big.Int).SetInt64(int64(mantissa))

	if exponent <= 3 {
		shift := 8 * (3 - exponent)
		target.Rsh(target, uint(shift))
	} else {
		shift := 8 * (exponent - 3)
		target.Lsh(target, uint(shift))
	}

	return target
}

func CompactToDifficulty(nBits uint32) *big.Float {
	target := CompactToTarget(nBits)

	diff1Target := CompactToTarget(0x1d00ffff)

	f1 := new(big.Float).SetInt(diff1Target)
	f2 := new(big.Float).SetInt(target)

	difficulty := new(big.Float).Quo(f1, f2)
	return difficulty
}
