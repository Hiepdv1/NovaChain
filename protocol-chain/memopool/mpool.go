package memopool

import (
	"bytes"
	blockchain "core-blockchain/core"
	"encoding/hex"
	"maps"
	"slices"
	"sort"

	log "github.com/sirupsen/logrus"
)

type TxInfo struct {
	Fee         float64
	Transaction blockchain.Transaction
}

type Memopool struct {
	Pending map[string]TxInfo
	Queued  map[string]TxInfo
}

func GetTxInfo(tx *blockchain.Transaction, bl *blockchain.Blockchain) *TxInfo {
	if !bl.VerifyTransaction(tx) {
		log.Infof("Transaction ID: %s is not valid", hex.EncodeToString(tx.ID))
		return nil
	}

	totalInput := 0.0
	for _, in := range tx.Inputs {
		prevTx, err := bl.FindTransaction(in.ID)
		if err != nil {
			log.Error("Transaction Not Found")
			continue
		}

		out := prevTx.Outputs[in.Out]
		totalInput += out.Value

	}

	totalOutput := 0.0
	for _, out := range tx.Outputs {
		totalOutput += out.Value
	}

	return &TxInfo{
		Fee:         totalInput - totalOutput,
		Transaction: *tx,
	}
}

func (memo *Memopool) GetTxByID(txID string) *blockchain.Transaction {
	if _, exists := memo.Pending[txID]; exists {
		tx := memo.Pending[txID].Transaction
		return &tx
	}

	if _, exists := memo.Queued[txID]; exists {
		tx := memo.Queued[txID].Transaction
		return &tx
	}

	return nil
}

func (memo *Memopool) Add(tx TxInfo) {
	txID := hex.EncodeToString(tx.Transaction.ID)
	if _, exists := memo.Pending[txID]; exists {
		return
	}

	if _, exists := memo.Queued[txID]; exists {
		return
	}

	memo.Pending[txID] = tx
}

func (memo *Memopool) HasPending(txID string) bool {
	_, exists := memo.Pending[txID]
	return exists
}

func (memo *Memopool) HashTX(txID string) bool {
	if _, exists := memo.Pending[txID]; exists {
		return true
	}

	if _, exists := memo.Queued[txID]; exists {
		return true
	}

	return false
}

func (memo *Memopool) GetTransactionHashes() (txs [][]byte) {
	for _, tx := range memo.Pending {
		txs = append(txs, tx.Transaction.ID)

	}
	return txs
}

func (memo *Memopool) RemoveFromAll(txID string) {
	delete(memo.Queued, txID)
	delete(memo.Pending, txID)
}

func (memo *Memopool) Move(tx TxInfo, to string) {
	if to == MEMO_MOVE_FLAG_PENDING {
		txID := hex.EncodeToString(tx.Transaction.ID)
		memo.Remove(txID, MEMO_MOVE_FLAG_QUEUED)
		memo.Pending[txID] = tx
	}

	if to == MEMO_MOVE_FLAG_QUEUED {
		txID := hex.EncodeToString(tx.Transaction.ID)
		memo.Remove(txID, MEMO_MOVE_FLAG_PENDING)
		memo.Queued[txID] = tx
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
	memo.Pending = map[string]TxInfo{}
	memo.Queued = map[string]TxInfo{}
}

func (memo *Memopool) SelectHighFeeTx() map[string]blockchain.Transaction {
	maxSizeBlock := blockchain.MaxBlockSize // mb

	// Reset queue before selecting highest-fee transactions for the next block
	memo.Queued = make(map[string]TxInfo, len(memo.Pending))

	totalSize := 0

	txPendings := slices.Collect(maps.Values(memo.Pending))

	sort.Slice(txPendings, func(i, j int) bool {
		return txPendings[i].Fee > txPendings[j].Fee
	})

	txs := make(map[string]blockchain.Transaction, 0)

	for _, tx := range txPendings {

		buf := new(bytes.Buffer)
		blockchain.SerializeTransaction(&tx.Transaction, buf)

		txSize := len(buf.Bytes())
		totalSize += txSize

		if totalSize > maxSizeBlock {
			break
		}

		memo.Move(tx, MEMO_MOVE_FLAG_QUEUED)

		txs[hex.EncodeToString(tx.Transaction.ID)] = tx.Transaction
	}

	return txs
}
