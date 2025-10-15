package utils

func FormatHashrate(hashrate float64) (float64, string) {
	switch {
	case hashrate >= 1e12:
		return hashrate / 1e12, "TH/s"
	case hashrate >= 1e9:
		return hashrate / 1e9, "GH/s"
	case hashrate >= 1e6:
		return hashrate / 1e6, "MH/s"
	case hashrate >= 1e3:
		return hashrate / 1e3, "kH/s"
	default:
		return hashrate, "H/s"
	}
}
