package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

type Blockchain struct {
	LastHash   []byte
	Database   *badger.DB
	InstanceId string
}

const (
	AdjustmentInterval  = 10
	TargetBlockTime     = 600 // 10 minute
	MaxTimestampDrift   = 600 // 10 minute
	MinDifficulty       = 24
	MaxDifficulty       = 34
	MaxDifficultyChange = 0.2
	CheckpointInterval  = 10
	BestHeightPrefix    = "lh"
	genesisData         = "genesis"
	CheckpointPrefix    = "checkpoint-"
)

var (
	// mutex         = &sync.Mutex{}
	_, file, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(file), "../")
)

func GetDatabasePath(port string) string {
	if port != "" {
		return filepath.Join(Root, fmt.Sprintf("./.chain/blocks_%s", port))
	}

	return filepath.Join(Root, "./.chain/blocks")
}

func DBExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func Exists(InstanceId string) bool {
	return DBExists(GetDatabasePath(InstanceId))
}

func InitBlockchain(instanceId string) *Blockchain {
	var lastHash []byte
	path := GetDatabasePath(instanceId)

	if DBExists(path) {
		log.Panic("Blockchain already exist")
		return nil
	} else {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic(err)
		}
	}

	opts := badger.DefaultOptions(path)
	opts.ValueDir = path
	db, err := OpenDB(path, opts)
	utils.ErrorHandle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := InitGenesisTx(genesisData)
		log.Info("No existing blockchain found")

		genesis := Genesis(cbtx)

		genesis.NChainWork = big.NewInt(0)

		genesis.NChainWork = genesis.NChainWork.Lsh(big.NewInt(1), uint(genesis.Difficulty))

		genesis.NChainWork = genesis.NChainWork.Div(big.NewInt(2).Lsh(big.NewInt(2), 256), genesis.NChainWork.Add(genesis.NChainWork, big.NewInt(1)))

		log.Info("NChainWork: ", genesis.NChainWork)

		err = txn.Set(genesis.Hash, genesis.Serialize())

		utils.ErrorHandle(err)
		err = txn.Set([]byte(BestHeightPrefix), genesis.Hash)
		lastHash = genesis.Hash

		key := fmt.Sprint(CheckpointPrefix, genesis.Height)
		err := txn.Set([]byte(key), genesis.Hash)

		return err
	})

	utils.ErrorHandle(err)

	chain := &Blockchain{lastHash, db, instanceId}

	utxo := UTXOSet{
		Blockchain: chain,
	}

	utxo.Compute()

	return chain
}

func OpenBadgerDB(instanceId string) *badger.DB {
	path := GetDatabasePath(instanceId)

	log.Info("Path: ", path)

	if !DBExists(path) {
		log.Info("No Existing Blockchian DB found, create one!")
		runtime.Goexit()
	}

	otps := badger.DefaultOptions(path).
		WithSyncWrites(true).
		WithTruncate(true)

	db, err := OpenDB(path, otps)

	utils.ErrorHandle(err)

	return db
}

func (bc *Blockchain) ContinueBlockchain() *Blockchain {
	var lastHash []byte
	var db *badger.DB

	if bc.Database == nil {
		db = OpenBadgerDB(bc.InstanceId)
	} else {
		db = bc.Database
	}

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(BestHeightPrefix))
		utils.ErrorHandle(err)

		lastHash, err = item.ValueCopy(nil)

		return err
	})

	if err != nil {
		lastHash = nil
	}

	return &Blockchain{LastHash: lastHash, Database: db, InstanceId: bc.InstanceId}

}

func (bc *Blockchain) HasBlock(hash []byte) (bool, error) {
	err := bc.Database.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(hash))
		return err
	})

	if err == badger.ErrKeyNotFound {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("error checking block in BadgerDB: %w", err)
	}

	return true, nil
}

func (bc *Blockchain) AdjustDifficulty(lastestBlock Block) int64 {

	if lastestBlock.Height%AdjustmentInterval != 0 {
		return int64(math.Max(MinDifficulty, float64(lastestBlock.Difficulty)))
	}

	blocks := make([]*Block, 0, AdjustmentInterval)

	currentBlockHash := lastestBlock.Hash

	for i := 0; i < AdjustmentInterval && len(currentBlockHash) > 0; i++ {
		block, err := bc.GetBlock(currentBlockHash)

		if err != nil {
			return int64(math.Max(MinDifficulty, float64(lastestBlock.Difficulty)))
		}

		blocks = append(blocks, &block)

		currentBlockHash = block.PrevHash
	}

	for i := 1; i < len(blocks); i++ {
		if blocks[i].Timestamp >= blocks[i-1].Timestamp {
			return int64(math.Max(MinDifficulty, float64(lastestBlock.Difficulty)))
		}

		if blocks[i-1].Timestamp-blocks[i].Timestamp >= 3600 {
			return int64(math.Max(MinDifficulty, float64(lastestBlock.Difficulty)))
		}
	}

	actualTime := int64(0)

	for i := 1; i < len(blocks); i++ {
		actualTime += blocks[i-1].Timestamp - blocks[i].Timestamp
	}

	actualTime /= int64(len(blocks))

	expectedTime := int64(TargetBlockTime)
	ratio := float64(actualTime) / float64(expectedTime)

	newDifficulty := float64(lastestBlock.Difficulty)

	if ratio < 0.5 {
		newDifficulty *= (1 + MaxDifficultyChange)
		newDifficulty = math.Ceil(newDifficulty)
	} else if ratio > 2 {
		newDifficulty *= (1 - MaxDifficultyChange)
		newDifficulty = math.Floor(newDifficulty)
	}

	newDifficulty = math.Max(float64(MinDifficulty), math.Min(MaxDifficulty, newDifficulty))

	log.Printf("Difficulty for block: height=%d, actualTime=%d, ratio=%.2f, newDifficulty=%d",
		lastestBlock.Height, actualTime, ratio, int(newDifficulty))

	return int64(newDifficulty)

}

func (bc *Blockchain) IsValidCheckpoint(bl *Block) bool {
	prevBlock, err := bc.GetBlock(bl.PrevHash)

	if err != nil {
		return false
	}

	checkpointHeight := (bl.Height / CheckpointInterval) * CheckpointInterval
	if checkpointHeight > 0 {
		var checkpointHash []byte
		err := bc.Database.View(func(txn *badger.Txn) error {
			key := fmt.Sprint(CheckpointPrefix, checkpointHeight)
			item, err := txn.Get([]byte(key))
			utils.ErrorHandle(err)

			checkpointHash, err = item.ValueCopy(nil)
			return err
		})

		utils.ErrorHandle(err)

		checkpointBlock, err := bc.GetBlock(checkpointHash)

		if err != nil || checkpointBlock.Height != checkpointHeight {
			log.Warn("Checkpoint validation failed")
			return false
		}

		currentHash := prevBlock.PrevHash
		currentHeight := prevBlock.Height

		for currentHeight > checkpointHeight && len(currentHash) > 0 {
			currentBlock, err := bc.GetBlock(currentHash)

			if err != nil {
				log.Warnf("Invalid chain: cannot find block with hash %x", currentHash)
				return false
			}

			if currentBlock.Height != currentHeight-1 {
				log.Warnf("Height mismatch at block %x: expected %d, got %d", currentHash, currentHeight-1, currentBlock.Height)
				return false
			}

			prevBlock, err := bc.GetBlock(currentBlock.PrevHash)

			if err != nil {
				log.Warnf("Missing previous block %x", currentBlock.PrevHash)
				return false
			}

			if currentBlock.Difficulty != bc.AdjustDifficulty(prevBlock) {
				log.Warnf("Difficulty mismatch at block %x", currentBlock.Hash)
				return false
			}

			currentHash = currentBlock.PrevHash
			currentHeight--

		}

		if checkpointBlock.Height != currentHeight {
			log.Warnf("Checkpoint mismatch: expected height %d, got %d", checkpointBlock.Height, currentHeight)
			return false
		}

		if !bytes.Equal(checkpointBlock.Hash, currentHash) {
			log.Warnf("Checkpoint mismatch: expected %x, got %x", checkpointBlock.Hash, currentHash)
			return false
		}

	}

	return true
}

func (bc *Blockchain) ComputeChain(lashBlockForkChain Block) error {
	block, err := bc.GetLastBlock()
	if err != nil {
		return err
	}
	lastHash := block.Hash

	currentBlock := lashBlockForkChain

	err = bc.Database.Update(func(txn *badger.Txn) error {
		for {

			err := txn.Delete(currentBlock.Hash)
			if err != nil {
				return err
			}

			if bytes.Equal(currentBlock.PrevHash, lastHash) {
				break
			}

			currentBlock, err = bc.GetBlock(currentBlock.PrevHash)

			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error computing chain: %w", err)
	}

	return nil
}

func (bc *Blockchain) IsBlockValid(bl Block) bool {
	prevBlock, err := bc.GetBlock(bl.PrevHash)
	utils.ErrorHandle(err)
	currentTime := time.Now().Unix()

	if bl.Difficulty != bc.AdjustDifficulty(prevBlock) {
		log.Warn("Difficulty validation failed")
		return false
	}

	if !bc.ValidateBlockTransactions(&bl) {
		log.Warn("Invalid Transaction")
		return false
	}

	if bl.Timestamp >= currentTime+MaxTimestampDrift || bl.Timestamp < prevBlock.Timestamp {
		log.Warnf("Invalid timestamp: too far in future or past. Current timestamp of block %d, MaxTimestampDrift %d", bl.Timestamp, currentTime+MaxDifficulty)
		return false
	}

	if !bc.IsValidCheckpoint(&bl) {
		return false
	}

	return bl.IsBlockValid(prevBlock)
}

func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {

	var block Block

	err := bc.Database.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("Block does not exist")
		} else {
			if blockData, err := item.ValueCopy(nil); err != nil {
				return err
			} else {
				block = *block.Deserialize(blockData)
				return nil
			}
		}
	})

	if err != nil {
		return block, err
	}

	return block, nil

}

func (bc *Blockchain) GetBlockHashes(blockHash []byte, height int64, max int64) [][]byte {
	var blocks [][]byte
	iter := bc.Iterator()

	if iter == nil {
		return blocks
	}

	for {

		block := iter.Next()

		if block.Height == height && bytes.Equal(block.Hash, blockHash) {
			break
		}

		blocks = append(blocks, block.PrevHash)

		if int64(len(blocks)) > max {
			blocks = blocks[int64(len(blocks))-max:]
		}

		if len(block.PrevHash) == 0 || block.PrevHash == nil {
			break
		}
	}

	return blocks
}

func (bc *Blockchain) GetLastBlock() (*Block, error) {
	var block Block

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(BestHeightPrefix))

		if err != nil {
			return err
		}

		hash, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		item, err = txn.Get(hash)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			block = *block.Deserialize(val)
			return nil
		})

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("get last block error: %s", err)
	}

	return &block, nil
}

func (bc *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXOs := make(map[string]TxOutputs)
	spentUTXOs := make(map[string][]int64)

	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentUTXOs[txID] != nil {
					for _, spentOut := range spentUTXOs[txID] {
						if spentOut == int64(outIdx) {
							continue Outputs
						}
					}
				}

				outs := UTXOs[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXOs[txID] = outs
			}

			if !tx.IsMinerTx() {
				for _, in := range tx.Inputs {
					inTxId := hex.EncodeToString(in.ID)
					spentUTXOs[inTxId] = append(spentUTXOs[inTxId], in.Out)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("No transaction with ID: " + string(ID))
}

func (bc *Blockchain) GetTransaction(transaction *Transaction) map[string]Transaction {
	txs := make(map[string]Transaction)

	for _, in := range transaction.Inputs {
		tx, err := bc.FindTransaction(in.ID)

		if err != nil {
			log.Error("Error: Invalid Transaction Ewww")
		}

		utils.ErrorHandle(err)

		txs[hex.EncodeToString(tx.ID)] = tx
	}

	return txs
}

func (bc *Blockchain) SignTransaction(privKey ecdsa.PrivateKey, tx *Transaction) {
	prevTxs := bc.GetTransaction(tx)
	tx.Sign(privKey, prevTxs)
}

func (bc *Blockchain) ValidateBlockTransactions(bl *Block) bool {
	utxos := bc.FindUTXO()

	for _, tx := range bl.Transactions {
		if tx.IsMinerTx() {
			continue
		}

		if !bc.VerifyTransaction(tx) {
			return false
		}

		for _, in := range tx.Inputs {
			txID := hex.EncodeToString(in.ID)
			if _, exists := utxos[txID]; !exists {
				return false
			}
		}
	}

	return true
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsMinerTx() {
		return true
	}

	prevTxs := bc.GetTransaction(tx)

	return tx.Verify(prevTxs)
}

func (bc *Blockchain) MineBlock(transactions []*Transaction, address string, callback func([]*Transaction)) (*Block, error) {

	publicKey := transactions[0].Inputs[0].ID
	var totalInput float64
	var totalOuput float64

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("Invalid Transaction")
		}

		for _, in := range tx.Inputs {
			if !bytes.Equal(publicKey, in.PubKey) {
				panic("Invalid transactions")
			}
			tx, err := bc.FindTransaction(in.ID)
			if err != nil {
				return nil, err
			}
			totalInput += tx.Outputs[in.Out].Value
		}

		for _, out := range tx.Outputs {
			totalOuput += out.Value
		}
	}

	lastestBlock, err := bc.GetLastBlock()
	if err != nil {
		return nil, err
	}

	difficulty := bc.AdjustDifficulty(*lastestBlock)

	transactions = append(transactions, bc.GetBlockReward(int64(lastestBlock.Height), address))

	block := CreateBlock(transactions, lastestBlock.Hash, lastestBlock.Height+1, difficulty)

	err = bc.AddBlock(block, callback)

	if err != nil {
		return nil, fmt.Errorf("error adding block: %v", err)
	}

	return block, nil
}

func (bc *Blockchain) GetBestHeight() int64 {
	var lastBlock Block

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(BestHeightPrefix))
		utils.ErrorHandle(err)

		lastHash, err := item.ValueCopy(nil)
		utils.ErrorHandle(err)

		item, err = txn.Get(lastHash)
		utils.ErrorHandle(err)

		lastBlockData, err := item.ValueCopy(nil)
		utils.ErrorHandle(err)

		lastBlock = *lastBlock.Deserialize(lastBlockData)

		return nil

	})

	if err != nil {
		utils.ErrorHandle(err)
	}

	return lastBlock.Height
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")

	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}

	retryOpts := originalOpts
	retryOpts.Truncate = true

	db, err := badger.Open(retryOpts)

	return db, err
}

func OpenDB(dir string, otps badger.Options) (*badger.DB, error) {

	otps.Logger = nil

	if db, err := badger.Open(otps); err != nil {

		if strings.Contains(err.Error(), "LOCK") {

			if db, err := retry(dir, otps); err == nil {
				lsm, vlog := db.Size()
				log.Infof("DB Size — LSM: %d bytes, ValueLog: %d bytes\n", lsm, vlog)
				log.Info("database unlocked , value log truncated")
				return db, nil
			}

			log.Panicln("could not unlock database: ", err)

		}

		return nil, err
	} else {
		return db, nil
	}

}
