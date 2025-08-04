package memopool

import (
	blockchain "core-blockchain/core"
	"encoding/hex"
)

type Memopool struct {
	Pending map[string]blockchain.Transaction
	Queued  map[string]blockchain.Transaction
}

func (memo *Memopool) Add(tnx blockchain.Transaction) {
	memo.Pending[hex.EncodeToString(tnx.ID)] = tnx
}

func (memo *Memopool) GetTransactions(count int64) (txs [][]byte) {
	var i int64 = 0
	for _, tx := range memo.Pending {
		txs = append(txs, tx.ID)
		i++
		if i == count {
			break
		}
	}
	return txs
}

func (memo *Memopool) RemoveFromAll(txID string) {
	delete(memo.Queued, txID)
	delete(memo.Pending, txID)
}

func (memo *Memopool) Move(tnx blockchain.Transaction, to string) {
	if to == MEMO_MOVE_FLAG_PENDING {
		txID := hex.EncodeToString(tnx.ID)
		memo.Remove(txID, MEMO_MOVE_FLAG_QUEUED)
		memo.Pending[txID] = tnx
	}

	if to == MEMO_MOVE_FLAG_QUEUED {
		txID := hex.EncodeToString(tnx.ID)
		memo.Remove(txID, MEMO_MOVE_FLAG_PENDING)
		memo.Queued[txID] = tnx
	}
}

func (memo *Memopool) Remove(txID string, from string) {
	if from == MEMO_MOVE_FLAG_QUEUED {
		delete(memo.Queued, txID)
		return
	}

	if from == MEMO_MOVE_FLAG_PENDING {
		delete(memo.Pending, txID)
		return
	}
}

func (memo *Memopool) ClearAll() {
	memo.Pending = map[string]blockchain.Transaction{}
	memo.Queued = map[string]blockchain.Transaction{}
}
