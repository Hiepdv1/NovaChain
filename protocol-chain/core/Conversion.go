package blockchain

import "math/big"

// Compact to Big
func CompactToBig(compact uint32) *big.Int {
	exponent := uint(compact >> 24)
	mantissa := compact & 0xFFFFFF
	target := new(big.Int).SetUint64(uint64(mantissa))
	if exponent <= 3 {
		target.Rsh(target, 8*(3-exponent))
	} else {
		target.Lsh(target, 8*(exponent-3))
	}
	return target
}

// Big to Compact
func BigToCompact(target *big.Int) uint32 {
	size := uint(len(target.Bytes()))
	var compact uint32
	if size <= 3 {
		compact = uint32(target.Uint64() << (8 * (3 - size)))
	} else {
		tmp := new(big.Int).Rsh(target, 8*(size-3))
		compact = uint32(tmp.Uint64())
	}
	if compact&0x00800000 != 0 {
		compact >>= 8
		size++
	}
	compact |= uint32(size) << 24
	return compact
}
