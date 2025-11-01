package utils

import (
	"math"
	"math/big"
)

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

	diff1Target := CompactToTarget(0x1f00ffff)

	f1 := new(big.Float).SetInt(diff1Target)
	f2 := new(big.Float).SetInt(target)

	difficulty := new(big.Float).Quo(f1, f2)
	return difficulty
}

func DifficultyFromNBits(nBits uint32) float64 {
	maxTarget := CompactToTarget(0x1f00ffff)
	target := CompactToTarget(nBits)

	ratio := new(big.Rat).SetFrac(maxTarget, target)
	f, _ := ratio.Float64()
	return f
}

func AverageDifficulty(nbitsList []uint32) float64 {
	if len(nbitsList) == 0 {
		return 0
	}
	var sum float64
	for _, n := range nbitsList {
		sum += DifficultyFromNBits(n)
	}
	return sum / float64(len(nbitsList))
}

func AverageBlockTime(timestamps []int64) float64 {
	if len(timestamps) < 2 {
		return 0
	}
	first := timestamps[len(timestamps)-1]
	last := timestamps[0]
	return float64(last-first) / float64(len(timestamps)-1)
}

func CalculateHashrate(avgDifficulty, avgBlockTime float64) float64 {
	if avgBlockTime == 0 {
		return 0
	}
	return avgDifficulty * math.Pow(2, 32) / avgBlockTime
}

func EvaluateNetwork(avgBlockTime, targetTime float64) string {
	if avgBlockTime <= 0 || targetTime <= 0 {
		return "Unknown"
	}

	ratio := avgBlockTime / targetTime

	switch {
	case ratio <= 0.5:
		return "Strong"
	case ratio <= 1.2:
		return "Normal"
	case ratio <= 2.0:
		return "Slow"
	default:
		return "Weak"
	}
}

func EvaluateSyncStatus(current, known int64) string {
	switch {
	case current < known-10:
		return "catching_up"
	case current < known:
		return "syncing"
	default:
		return "synced"
	}
}

func CalculateHashrateChange(currentHashrate, previousHashrate float64) float64 {
	changePercent := ((currentHashrate - previousHashrate) / previousHashrate) * 100
	return math.Round(changePercent*100) / 100
}
