package constants

type PriorityInfo struct {
	Name string
	Rate float64
}

type TxStatus string

const (
	TxStatusPending TxStatus = "pending" // transaction in mempool, waiting to be picked
	TxStatusMined   TxStatus = "mined"   // included in a block, not yet finalized

	TxStatusFailed TxStatus = "failed" // invalid transaction (e.g., insufficient funds, bad signature)
)

const (
	PriorityLow    = 1
	PriorityNormal = 2
	PriorityHigh   = 3
)

var Priorities = map[uint]PriorityInfo{
	PriorityLow: {
		Name: "Low",
		Rate: 0.1,
	},
	PriorityNormal: {
		Name: "Normal",
		Rate: 0.6,
	},
	PriorityHigh: {
		Name: "High",
		Rate: 0.9,
	},
}
