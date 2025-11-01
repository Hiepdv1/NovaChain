package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

var (
	mutex = &sync.Mutex{}
)

const (
	ChainPrefix   = "chain-"
	MaxForkLength = 6
)

func (bc *Blockchain) Reorganize(newBlock *Block, callback func([]*Transaction)) error {
	log.Infof("Reorg started: candidate block %x at height %d", newBlock.Hash, newBlock.Height)

	newChain := make([]*Block, 0)
	newHash := newBlock.PrevHash
	newHeight := newBlock.Height - 1
	maxForklength := MaxForkLength

	for len(newHash) > 0 && newHeight > 0 {
		block, err := bc.GetBlock(newHash)
		if err != nil {
			return fmt.Errorf("failed to get block %x in new chain: %w", newHash, err)
		}
		newChain = append(newChain, &block)
		newHash = block.PrevHash
		newHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}
	log.Infof("Reorg: collected %d blocks in new chain (maxFork=%d)", len(newChain), MaxForkLength)

	currentTip, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("failed to get last block: %w", err)
	}
	log.Infof("Reorg: current tip %x at height %d", currentTip.Hash, currentTip.Height)

	currentHash := currentTip.Hash
	currentHeight := currentTip.Height
	maxForklength = MaxForkLength
	oldChain := make([]*Block, 0)

	for len(currentHash) > 0 && currentHeight > 0 {
		block, err := bc.GetBlock(currentHash)
		if err != nil {
			return fmt.Errorf("failed to get block %x in old chain: %w", currentHash, err)
		}
		oldChain = append(oldChain, &block)
		currentHash = block.PrevHash
		currentHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}
	log.Infof("Reorg: collected %d blocks in old chain (maxFork=%d)", len(oldChain), MaxForkLength)

	var commonAncestorHeight int64 = 0
	var commonAncestorHash []byte
	for _, newBlock := range newChain {
		for _, oldBlock := range oldChain {
			if newBlock.Height == oldBlock.Height && bytes.Equal(newBlock.Hash, oldBlock.Hash) {
				commonAncestorHeight = newBlock.Height
				commonAncestorHash = newBlock.Hash
				log.Infof("Reorg: found common ancestor at height %d hash=%x", commonAncestorHeight, commonAncestorHash)
				break
			}
		}
		if len(commonAncestorHash) > 0 {
			break
		}
	}

	if commonAncestorHeight == 0 {
		return errors.New("reorg failed: no common ancestor found or fork too deep")
	}

	var memoryPool []*Transaction
	UTXOs, err := bc.FindUTXO()
	if err != nil {
		return err
	}

	for _, oldBlock := range oldChain {
		if oldBlock.Height < commonAncestorHeight {
			continue
		}
		for _, tx := range oldBlock.Transactions {
			if tx.IsMinerTx() {
				continue
			}
			txID := hex.EncodeToString(tx.ID)
			if _, exists := UTXOs[txID]; !exists {
				memoryPool = append(memoryPool, tx)
				log.Infof("Reorg: rollback tx %s to mempool", txID)
			}
		}
	}

	for _, nb := range newChain {
		if nb.Height < commonAncestorHeight {
			continue
		}
		for _, tx := range nb.Transactions {
			for ix, memoTx := range memoryPool {
				if bytes.Equal(tx.ID, memoTx.ID) {
					memoryPool = append(memoryPool[:ix], memoryPool[ix+1:]...)
					log.Infof("Reorg: tx %x already in new chain, removed from mempool rollback list", tx.ID)
				}
			}
		}
	}

	if len(memoryPool) > 0 {
		callback(memoryPool)
		txIds := make([]string, 0)
		for _, tx := range memoryPool {
			txIds = append(txIds, hex.EncodeToString(tx.ID))
		}
		log.Warnf("Reorg: rollback %d txs from old chain:\n%s", len(memoryPool), strings.Join(txIds, "\n"))
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		log.Warnf("Reorg: switching to new best chain, new tip %x at height %d", newBlock.Hash, newBlock.Height)

		if err := txn.Set([]byte(BestHeightPrefix), newBlock.Hash); err != nil {
			return fmt.Errorf("failed to set best height: %w", err)
		}
		bc.LastHash = newBlock.Hash

		keyCheckpoint := fmt.Sprintf("%s%d", CheckpointPrefix, newBlock.Height)
		if err := txn.Set([]byte(keyCheckpoint), newBlock.Hash); err != nil {
			return fmt.Errorf("failed to set checkpoint: %w", err)
		}

		serialize := SerializeBlock(newBlock)

		if err := txn.Set(newBlock.Hash, serialize); err != nil {
			return fmt.Errorf("failed to persist new block: %w", err)
		}
		return nil
	})

	if err == nil {
		bc.LastHash = newBlock.Hash
		log.Infof("Reorg finished successfully: new tip %x height=%d", newBlock.Hash, newBlock.Height)
	}

	return err
}

func (bc *Blockchain) CalcWork(nBits uint32) *big.Int {
	target := CompactToBig(nBits)
	denominator := new(big.Int).Add(target, big.NewInt(1))
	return new(big.Int).Div(
		new(big.Int).Lsh(big.NewInt(1), 256),
		denominator,
	)
}

func (bc *Blockchain) AddBlock(block *Block, callback func([]*Transaction)) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Infof("ADD BLOCK: processing block %x (height=%d, prev=%x)", block.Hash, block.Height, block.PrevHash)

	if !bc.IsBlockValid(*block) {
		return fmt.Errorf("ADD BLOCK: invalid block %x", block.Hash)
	}

	currentTipBlock, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("ADD BLOCK: failed to get current tip: %w", err)
	}
	log.Debugf("ADD BLOCK: current tip block is %x (height=%d)", currentTipBlock.Hash, currentTipBlock.Height)

	prevBlock, err := bc.GetBlock(block.PrevHash)
	if err != nil {
		return fmt.Errorf("ADD BLOCK: failed to get previous block %x: %w", block.PrevHash, err)
	}
	log.Debugf("ADD BLOCK: previous block %x loaded successfully", block.PrevHash)

	newChainWork := bc.CalcWork(block.NBits)
	newChainWork = new(big.Int).Add(prevBlock.NChainWork, newChainWork)
	block.NChainWork = newChainWork

	log.Infof("ADD BLOCK: received block %x (height=%d). Current tip work=%s → new block chain work=%s",
		block.Hash, block.Height, currentTipBlock.NChainWork.String(), newChainWork.String(),
	)

	if newChainWork.Cmp(currentTipBlock.NChainWork) > 0 {
		log.Infof("ADD BLOCK: block %x has higher chain work → triggering chain reorganization", block.Hash)

		if err := bc.Reorganize(block, callback); err != nil {
			return fmt.Errorf("ADD BLOCK: reorganization failed for block %x: %w", block.Hash, err)
		}

		log.Infof("ADD BLOCK: reorganization completed successfully for block %x", block.Hash)
	} else {
		log.Infof("ADD BLOCK: block %x has lower chain work → saving as disconnected block", block.Hash)

		err := bc.Database.Update(func(txn *badger.Txn) error {
			keyCheckpoint := fmt.Sprint(CheckpointPrefix, block.Height)
			if err := txn.Set([]byte(keyCheckpoint), block.Hash); err != nil {
				return err
			}
			serialize := SerializeBlock(block)
			return txn.Set(block.Hash, serialize)
		})

		if err != nil {
			return fmt.Errorf("ADD BLOCK: failed to save disconnected block %x: %v", block.Hash, err)
		}

		log.Debugf("ADD BLOCK: disconnected block %x saved successfully", block.Hash)
	}

	log.Info("ADD BLOCK: updating UTXO set...")
	utxoSet := UTXOSet{
		Blockchain: bc,
	}

	err = utxoSet.Compute()
	if err != nil {
		return fmt.Errorf("ADD BLOCK: failed to compute UTXO set: %v", err)
	}

	log.Info("ADD BLOCK: block successfully added and UTXO set updated")
	return nil
}
