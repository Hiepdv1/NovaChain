package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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
	log.Infof("🔄 Reorg started — candidate block=%x height=%d", newBlock.Hash[:6], newBlock.Height)

	newChain := make([]*Block, 0)
	newHash := newBlock.PrevHash
	newHeight := newBlock.Height - 1
	maxForklength := MaxForkLength

	for len(newHash) > 0 && newHeight > 0 {
		block, err := bc.GetBlock(newHash)
		if err != nil {
			return fmt.Errorf("❌ Reorg: get block %x in new chain: %w", newHash[:6], err)
		}
		newChain = append(newChain, &block)
		newHash = block.PrevHash
		newHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}
	log.Infof("📦 Reorg: collected %d new blocks (maxFork=%d)", len(newChain), MaxForkLength)

	currentTip, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("❌ Reorg: get current tip: %w", err)
	}
	log.Infof("⛓️ Current tip %x height=%d", currentTip.Hash[:6], currentTip.Height)

	currentHash := currentTip.Hash
	currentHeight := currentTip.Height
	maxForklength = MaxForkLength
	oldChain := make([]*Block, 0)

	for len(currentHash) > 0 && currentHeight > 0 {
		block, err := bc.GetBlock(currentHash)
		if err != nil {
			return fmt.Errorf("❌ Reorg: get block %x in old chain: %w", currentHash[:6], err)
		}
		oldChain = append(oldChain, &block)
		currentHash = block.PrevHash
		currentHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}
	log.Infof("📦 Reorg: collected %d old blocks (maxFork=%d)", len(oldChain), MaxForkLength)

	var commonAncestorHeight int64
	var commonAncestorHash []byte
	for _, nb := range newChain {
		for _, ob := range oldChain {
			if nb.Height == ob.Height && bytes.Equal(nb.Hash, ob.Hash) {
				commonAncestorHeight = nb.Height
				commonAncestorHash = nb.Hash
				log.Infof("🔗 Found common ancestor — height=%d hash=%x", nb.Height, nb.Hash[:6])
				break
			}
		}
		if len(commonAncestorHash) > 0 {
			break
		}
	}

	if commonAncestorHeight == 0 {
		return errors.New("⚠️ Reorg failed — no common ancestor or fork too deep")
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
				log.Debugf("↩️ Rollback tx %s to mempool", txID[:8])
			}
		}
	}

	for _, nb := range newChain {
		if nb.Height < commonAncestorHeight {
			continue
		}

		for ix, memoTx := range memoryPool {
			for _, tx := range nb.Transactions {
				if bytes.Equal(tx.ID, memoTx.ID) {
					memoryPool = append(memoryPool[:ix], memoryPool[ix+1:]...)
					log.Debugf("🧹 Removed tx %x from rollback list (already in new chain)", tx.ID[:6])
				}
			}
		}

		err := bc.Database.Update(func(txn *badger.Txn) error {
			keyCheckpoint := fmt.Sprintf("%s%d", CheckpointPrefix, nb.Height)
			return txn.Set([]byte(keyCheckpoint), nb.Hash)
		})
		if err != nil {
			return err
		}
	}

	if len(memoryPool) > 0 {
		callback(memoryPool)
		log.Warnf("⚠️ Reorg rollback %d tx(s) from old chain", len(memoryPool))
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		log.Infof("🏁 Switching to new best chain — tip=%x height=%d", newBlock.Hash[:6], newBlock.Height)
		if err := txn.Set([]byte(BestHeightPrefix), newBlock.Hash); err != nil {
			return err
		}
		bc.LastHash = newBlock.Hash

		keyCheckpoint := fmt.Sprintf("%s%d", CheckpointPrefix, newBlock.Height)
		if err := txn.Set([]byte(keyCheckpoint), newBlock.Hash); err != nil {
			return err
		}

		return txn.Set(newBlock.Hash, SerializeBlock(newBlock))
	})

	if err == nil {
		log.Infof("✅ Reorg completed — new tip=%x height=%d", newBlock.Hash[:6], newBlock.Height)
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

	log.Infof("📥 Add block — %x height=%d prev=%x", block.Hash[:6], block.Height, block.PrevHash[:6])

	if !bc.IsBlockValid(*block) {
		return fmt.Errorf("❌ Invalid block %x", block.Hash[:6])
	}

	currentTip, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("❌ Failed to get current tip: %w", err)
	}
	prevBlock, err := bc.GetBlock(block.PrevHash)
	if err != nil {
		return fmt.Errorf("❌ Failed to get previous block %x: %w", block.PrevHash[:6], err)
	}

	newChainWork := bc.CalcWork(block.NBits)
	newChainWork = new(big.Int).Add(prevBlock.NChainWork, newChainWork)
	block.NChainWork = newChainWork

	log.Infof("⚙️ Chain work: tip=%s → new=%s", currentTip.NChainWork.String(), newChainWork.String())

	if newChainWork.Cmp(currentTip.NChainWork) > 0 {
		if bytes.Equal(currentTip.Hash, block.PrevHash) {
			err := bc.Database.Update(func(txn *badger.Txn) error {
				if err := txn.Set([]byte(BestHeightPrefix), block.Hash); err != nil {
					return err
				}
				bc.LastHash = block.Hash

				keyCheckpoint := fmt.Sprintf("%s%d", CheckpointPrefix, block.Height)
				if err := txn.Set([]byte(keyCheckpoint), block.Hash); err != nil {
					return err
				}

				return txn.Set(block.Hash, SerializeBlock(block))
			})

			if err != nil {
				return fmt.Errorf("❌ Add block failed: %x %w", block.Hash[:6], err)
			}

		} else {
			log.Infof("🔄 Reorg needed — block=%x", block.Hash[:6])
			if err := bc.Reorganize(block, callback); err != nil {
				return fmt.Errorf("❌ Reorg failed for block %x: %w", block.Hash[:6], err)
			}
			log.Infof("✅ Reorg done for block %x", block.Hash[:6])
		}
	} else {
		log.Warnf("🧩 Lower chain work — saving disconnected block %x", block.Hash[:6])
		err := bc.Database.Update(func(txn *badger.Txn) error {
			return txn.Set(block.Hash, SerializeBlock(block))
		})
		if err != nil {
			return fmt.Errorf("❌ Save disconnected block %x: %v", block.Hash[:6], err)
		}
		log.Infof("💾 Disconnected block %x saved", block.Hash[:6])
	}

	log.Info("🔢 Updating UTXO set...")
	utxoSet := UTXOSet{Blockchain: bc}
	if err := utxoSet.Compute(); err != nil {
		return fmt.Errorf("❌ UTXO update failed: %v", err)
	}

	log.Infof("✅ Block %x added successfully", block.Hash[:6])
	return nil
}
