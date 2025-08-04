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
	mutex.Lock()
	defer mutex.Unlock()

	currentTip, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("failed to get last block: %w", err)
	}

	newChain := make([]*Block, 0)
	newHash := newBlock.PrevHash
	newHeight := newBlock.Height - 1
	maxForklength := MaxForkLength

	for len(newHash) > 0 && newHeight > 0 {
		block, err := bc.GetBlock(newHash)
		if err != nil {
			return (fmt.Errorf("failed to get block %x in new chain: %w", newHash, err))
		}
		newChain = append(newChain, &block)
		newHash = block.PrevHash
		newHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}

	currentHash := currentTip.Hash
	currentHeight := currentTip.Height
	maxForklength = MaxForkLength
	oldChain := make([]*Block, 0)

	for len(currentHash) > 0 && currentHeight > 0 {
		block, err := bc.GetBlock(currentHash)
		if err != nil {
			return (fmt.Errorf("failed to get block %x in old chain: %w", currentHash, err))
		}
		oldChain = append(oldChain, &block)
		currentHash = block.PrevHash
		currentHeight--
		maxForklength--
		if maxForklength == 0 {
			break
		}
	}

	var forkHeight int64 = 0
	var forkHash []byte
	for _, newBlock := range newChain {
		for _, oldBlock := range oldChain {
			if newBlock.Height == oldBlock.Height && bytes.Equal(newBlock.Hash, oldBlock.Hash) {
				forkHeight = newBlock.Height
				forkHash = newBlock.Hash
			}
		}

		if len(forkHash) > 0 {
			break
		}
	}

	if forkHeight == 0 && len(forkHash) == 0 {
		return errors.New("no common ancestor found or Reorg too deep")
	}

	var memoryPool []*Transaction
	UTXOs := bc.FindUTXO()
	for _, oldBlock := range oldChain {
		if oldBlock.Height < forkHeight {
			continue
		}
		for _, tx := range oldBlock.Transactions {
			if tx.IsMinerTx() {
				continue
			}
			txID := hex.EncodeToString(tx.ID)
			if _, exists := UTXOs[txID]; !exists {
				memoryPool = append(memoryPool, tx)
				log.Infof("Reorg: Transaction %s returned to memopool", txID)
			}
		}
	}

	for _, newBlock := range newChain {
		if newBlock.Height < forkHeight {
			continue
		}

		for _, tx := range newBlock.Transactions {
			for ix, memoTx := range memoryPool {
				if bytes.Equal(tx.ID, memoTx.ID) {
					memoryPool = append(memoryPool[:ix], memoryPool[ix+1:]...)
				}
			}
		}
	}

	go callback(memoryPool)

	err = bc.Database.Update(func(txn *badger.Txn) error {
		log.Warningf("REORGNIZATION: New best chain found ! Tips: %x Height: %d", newBlock.Hash, newBlock.Height)

		err := txn.Set([]byte(BestHeightPrefix), newBlock.Hash)
		if err != nil {
			return fmt.Errorf("failed to set best height: %w", err)
		}
		bc.LastHash = newBlock.Hash

		keyCheckpoint := fmt.Sprintf("%s%d", CheckpointPrefix, newBlock.Height)
		err = txn.Set([]byte(keyCheckpoint), newBlock.Hash)
		if err != nil {
			return fmt.Errorf("failed to set best height: %w", err)
		}

		err = txn.Set(newBlock.Hash, newBlock.Serialize())

		return err
	})

	return err
}

func (bc *Blockchain) GetTotalWork(block *Block) (*big.Int, error) {

	prevBlock, err := bc.GetBlock(block.PrevHash)

	if err != nil {
		return nil, fmt.Errorf("failed to get block %x during total work calculation: %w", block.Hash, err)
	}

	if !bc.IsBlockValid(prevBlock) {
		return nil, fmt.Errorf("block %x not valid", block.Hash)
	}

	totalWork := big.NewInt(0)
	totalWork = totalWork.Add(totalWork, block.NChainWork)

	currentWork := big.NewInt(0)
	currentWork = currentWork.Lsh(big.NewInt(1), uint(block.Difficulty))

	currentWork = currentWork.Div(big.NewInt(2).Lsh(big.NewInt(2), 256), currentWork.Add(currentWork, big.NewInt(1)))

	totalWork = totalWork.Add(totalWork, currentWork)

	return totalWork, nil

}

func (bc *Blockchain) AddBlock(block *Block, callback func([]*Transaction)) error {
	mutex.Lock()
	defer mutex.Unlock()
	if !bc.IsBlockValid(*block) {
		return fmt.Errorf("block %x not valid", block.Hash)
	}

	currentTipBlock, err := bc.GetLastBlock()
	if err != nil {
		return fmt.Errorf("failed to get last block: %w", err)
	}

	mainChainWork, err := bc.GetTotalWork(currentTipBlock)

	if err != nil {
		return fmt.Errorf("failed to calculate work for main chain: %v", err)
	}

	newChainWork, err := bc.GetTotalWork(block)

	if err != nil {
		return fmt.Errorf("failed to calculate work for new chain: %v", err)
	}

	block.NChainWork = newChainWork

	log.Infof("Comparing chains: Main chain work = %s, New chain work = %s", mainChainWork.String(), newChainWork.String())
	if newChainWork.Cmp(mainChainWork) > 0 {
		err := bc.Reorganize(block, callback)
		if err != nil {
			return err
		}
	} else {
		err := bc.Database.Update(func(txn *badger.Txn) error {
			log.Info("UPDATE NEW BLOCK DISCONNECTED")

			keyCheckpoint := fmt.Sprint(CheckpointPrefix, block.Height)
			err = txn.Set([]byte(keyCheckpoint), block.Hash)
			if err != nil {
				return err
			}

			err = txn.Set(block.Hash, block.Serialize())

			return err
		})

		if err != nil {
			return fmt.Errorf("failed to update database with new block %x: %v", block.Hash, err)
		}
	}

	return nil
}
