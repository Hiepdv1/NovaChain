package dashboard

import dbchain "ChainServer/internal/db/chain"

type NetworkOverview struct {
	Chain struct {
		BestHeight int64
		PerHours   int64
	}
	Hashrate struct {
		Value  string
		Per24H string
	}
	Transaction struct {
		Total      int64
		AddedToday int64
	}
	PendingTx struct {
		Count      int64
		AddedToday int64
	}
	ActiveMiners struct {
		Count  int64
		Worker int64
	}
}

type RecentActivity struct {
	Blocks []dbchain.Block
	Txs    []dbchain.Transaction
}
