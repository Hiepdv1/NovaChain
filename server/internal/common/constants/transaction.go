package constants

type PriorityInfo struct {
	Name string
	Rate float64
}

type TxStatus string

const (
	TxStatusPending   TxStatus = "pending"   // transaction in mempool, waiting to be picked
	TxStatusMining    TxStatus = "mining"    // selected by miner, in candidate block
	TxStatusMined     TxStatus = "mined"     // included in a block, not yet finalized
	TxStatusConfirmed TxStatus = "confirmed" // block finalized, transaction safe

	TxStatusFailed  TxStatus = "failed"  // invalid transaction (e.g., insufficient funds, bad signature)
	TxStatusDropped TxStatus = "dropped" // removed from mempool (e.g., timeout, replaced)
	TxStatusReorged TxStatus = "reorged" // transaction was in block, but removed due to chain reorg
)

const (
	PriorityLow    = 0
	PriorityNormal = 1
	PriorityHigh   = 2
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
