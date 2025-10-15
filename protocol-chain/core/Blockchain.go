package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
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

var (
	MaxTarget = big.NewInt(0x1d00ffff)
)

const (
	AdjustmentInterval = 10
	TargetBlockTime    = 600 // 10 minute
	MaxTimestampDrift  = 600 // 10 minute
	CheckpointInterval = 10
	MaxBlockSize       = 1 * 1024 * 1024 // 1mb
	BestHeightPrefix   = "lh"
	CheckpointPrefix   = "checkpoint-"
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

func InitBlockchain(instanceId string) (*Blockchain, error) {
	var lastHash []byte
	path := GetDatabasePath(instanceId)

	if DBExists(path) {
		return nil, fmt.Errorf("%s", "Blockchain already exist")
	} else {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	opts := badger.DefaultOptions(path)
	opts.ValueDir = path
	db, err := OpenDB(path, opts)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(txn *badger.Txn) error {
		cbtx, err := InitGenesisTx(1)
		if err != nil {
			return err
		}
		log.Info("No existing blockchain found")

		genesis, err := Genesis(cbtx)
		if err != nil {
			return err
		}

		target := CompactToBig(genesis.NBits)

		denominator := new(big.Int).Add(target, big.NewInt(1))

		work := new(big.Int).Div(
			new(big.Int).Lsh(big.NewInt(1), 256),
			denominator)

		genesis.NChainWork = work

		log.Infof("NChainWork: %d", genesis.NChainWork)

		serialize := SerializeBlock(genesis)

		err = txn.Set(genesis.Hash, serialize)

		if err != nil {
			return err
		}

		err = txn.Set([]byte(BestHeightPrefix), genesis.Hash)

		if err != nil {
			return err
		}
		lastHash = genesis.Hash

		key := fmt.Sprint(CheckpointPrefix, genesis.Height)
		err = txn.Set([]byte(key), genesis.Hash)

		return err
	})

	if err != nil {
		return nil, err
	}

	chain := &Blockchain{lastHash, db, instanceId}

	utxo := UTXOSet{
		Blockchain: chain,
	}

	utxo.Compute()

	return chain, nil
}

func OpenBadgerDB(instanceId string) (*badger.DB, error) {
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

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (bc *Blockchain) ContinueBlockchain() (*Blockchain, error) {
	var lastHash []byte
	var db *badger.DB

	if bc.Database == nil {
		database, err := OpenBadgerDB(bc.InstanceId)
		if err != nil {
			return nil, err
		}
		db = database
	} else {
		db = bc.Database
	}

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(BestHeightPrefix))
		if err != nil {
			return err
		}

		lastHash, err = item.ValueCopy(nil)

		return err
	})

	if err != nil {
		lastHash = nil
	}

	return &Blockchain{LastHash: lastHash, Database: db, InstanceId: bc.InstanceId}, nil

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

func (bc *Blockchain) AdjustDifficulty(lastBlock *Block) uint32 {

	if lastBlock.Height%AdjustmentInterval != 0 {
		return lastBlock.NBits
	}

	firstHeight := lastBlock.Height - AdjustmentInterval
	firstBlock, err := bc.GetBlockByHeight(firstHeight)
	if err != nil {
		log.Error(err)
		return lastBlock.NBits
	}

	actualTimespan := lastBlock.Timestamp - firstBlock.Timestamp
	targetTimespan := int64(TargetBlockTime * AdjustmentInterval)

	if actualTimespan < targetTimespan/4 {
		actualTimespan = targetTimespan / 4
	}
	if actualTimespan > targetTimespan*4 {
		actualTimespan = targetTimespan * 4
	}

	oldTarget := CompactToBig(lastBlock.NBits)
	newTarget := new(big.Int).Mul(oldTarget, big.NewInt(actualTimespan))
	newTarget.Div(newTarget, big.NewInt(targetTimespan))

	if newTarget.Cmp(MaxTarget) > 0 {
		newTarget.Set(MaxTarget)
	}

	return BigToCompact(newTarget)
}

func (bc *Blockchain) IsValidCheckpoint(bl *Block) (bool, error) {
	prevBlock, err := bc.GetBlock(bl.PrevHash)

	if err != nil {
		return false, err
	}

	checkpointHeight := (bl.Height / CheckpointInterval) * CheckpointInterval
	if checkpointHeight > 0 {
		var checkpointHash []byte
		err := bc.Database.View(func(txn *badger.Txn) error {
			key := fmt.Sprint(CheckpointPrefix, checkpointHeight)
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			checkpointHash, err = item.ValueCopy(nil)
			return err
		})

		if err != nil {
			return false, err
		}

		checkpointBlock, err := bc.GetBlock(checkpointHash)

		if err != nil || checkpointBlock.Height != checkpointHeight {
			return false, err
		}

		currentHash := prevBlock.PrevHash
		currentHeight := prevBlock.Height

		for currentHeight > checkpointHeight && len(currentHash) > 0 {
			currentBlock, err := bc.GetBlock(currentHash)

			if err != nil {
				log.Warnf("Invalid chain: cannot find block with hash %x", currentHash)
				return false, nil
			}

			if currentBlock.Height != currentHeight-1 {
				log.Warnf("Height mismatch at block %x: expected %d, got %d", currentHash, currentHeight-1, currentBlock.Height)
				return false, nil
			}

			_, err = bc.GetBlock(currentBlock.PrevHash)

			if err != nil {
				log.Warnf("Missing previous block %x", currentBlock.PrevHash)
				return false, nil
			}

			currentHash = currentBlock.PrevHash
			currentHeight--

		}

		if checkpointBlock.Height != currentHeight {
			log.Warnf("Checkpoint mismatch: expected height %d, got %d", checkpointBlock.Height, currentHeight)
			return false, nil
		}

		if !bytes.Equal(checkpointBlock.Hash, currentHash) {
			log.Warnf("Checkpoint mismatch: expected %x, got %x", checkpointBlock.Hash, currentHash)
			return false, nil
		}

	}

	return true, nil
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
	if err != nil {
		log.Error(err)
		return false
	}
	currentTime := time.Now().Unix()

	if !bc.ValidateBlockTransactions(&bl) {
		log.Warn("Invalid Transaction")
		return false
	}

	if bl.Timestamp >= currentTime+MaxTimestampDrift || bl.Timestamp < prevBlock.Timestamp {
		log.Warnf("Invalid timestamp: too far in future or past. Current timestamp of block %d, MaxTimestampDrift %d", bl.Timestamp, currentTime+MaxTimestampDrift)
		return false
	}

	isValidCheckpoint, err := bc.IsValidCheckpoint(&bl)
	if err != nil {
		log.Error(err)
		return false
	}

	if !isValidCheckpoint {
		return false
	}

	if bl.NBits != bc.AdjustDifficulty(&prevBlock) {
		return false
	}

	size, err := bl.Size()
	if err != nil {
		log.Error(err)
		return false
	}

	if size > MaxBlockSize {
		return false
	}

	return bl.IsBlockValid(prevBlock)
}

func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {

	var block Block

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(blockHash)

		if err != nil {
			return err
		}

		blockData, err := item.ValueCopy(nil)

		if err != nil {
			return err
		}

		b := DeserializeBlockData(blockData)

		block = *b

		return nil
	})

	if err != nil {
		return Block{}, err
	}

	return block, nil

}

func (bc *Blockchain) GetBlockByHeight(height int64) (*Block, error) {
	hash, err := bc.GetBlockHashByHeight(height)
	if err != nil {
		return nil, err
	}

	block, err := bc.GetBlock(hash)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (bc *Blockchain) GetBlockHashByHeight(height int64) ([]byte, error) {
	var blockHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s%d", CheckpointPrefix, height)

		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		blockHash = value

		return nil
	})

	if err != nil {
		return nil, err
	}

	return blockHash, nil
}

func (bc *Blockchain) GetBlockLocator() ([][]byte, error) {

	var locator [][]byte
	step := 1
	height, err := bc.GetBestHeight()
	if err != nil {
		return nil, err
	}

	for height > 1 {
		hash, err := bc.GetBlockHashByHeight(height)
		if err != nil {
			return nil, err
		}

		locator = append(locator, hash)

		if len(locator) < 10 {
			height--
		} else {
			height -= int64(step)
			step *= 2
		}

	}

	genesisHash, err := bc.GetBlockHashByHeight(1)
	if err != nil {
		return nil, err
	}

	locator = append(locator, genesisHash)

	return locator, nil
}

func (bc *Blockchain) GetBlockHashes(blockHash []byte, height int64, max int64) ([][]byte, error) {
	var blocks [][]byte
	iter, err := bc.Iterator()

	if err != nil {
		return nil, err
	}

	for {

		block, err := iter.Next()
		if err != nil {
			return nil, err
		}

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

	return blocks, nil
}

func (bc *Blockchain) GetBlockRange(blockHash []byte, max int64) ([]Block, error) {
	var blocks []Block
	iter, err := bc.Iterator()

	if err != nil {
		return nil, err
	}

	for {

		block, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if len(block.PrevHash) == 0 || block.PrevHash == nil || bytes.Equal(block.Hash, blockHash) {
			break
		}

		blocks = append(blocks, *block)

		if int64(len(blocks)) > max {
			blocks = blocks[int64(len(blocks))-max:]
		}
	}

	return blocks, nil
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
			b := DeserializeBlockData(val)

			block = *b

			return nil
		})

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("get last block error: %s", err)
	}

	return &block, nil
}

func (bc *Blockchain) FindUTXO() (map[string]TxOutputs, error) {
	UTXOs := make(map[string]TxOutputs)
	spentUTXOs := make(map[string][]int64)

	iter, err := bc.Iterator()

	if err != nil {
		return nil, err
	}

	for {
		block, err := iter.Next()
		if err != nil {
			return nil, err
		}

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

	return UTXOs, nil
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter, err := bc.Iterator()
	if err != nil {
		return Transaction{}, nil
	}

	for {
		block, err := iter.Next()

		if err != nil {
			return Transaction{}, err
		}

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("No transaction with ID: " + hex.EncodeToString(ID))
}

func (bc *Blockchain) GetTransaction(transaction *Transaction) map[string]Transaction {
	txs := make(map[string]Transaction)

	for _, in := range transaction.Inputs {
		tx, err := bc.FindTransaction(in.ID)

		if err != nil {
			log.Errorf("Error: Find Transaction With Error %v", err)
			return nil
		}

		txs[hex.EncodeToString(tx.ID)] = tx
	}

	return txs
}

func (bc *Blockchain) SignTransaction(privKey ecdsa.PrivateKey, tx *Transaction) {
	prevTxs := bc.GetTransaction(tx)
	tx.Sign(privKey, prevTxs)
}

func (bc *Blockchain) ValidateBlockTransactions(bl *Block) bool {
	utxos, err := bc.FindUTXO()
	if err != nil {
		return false
	}

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

	utxoSet := UTXOSet{
		Blockchain: bc,
	}

	for _, in := range tx.Inputs {
		outs, _, err := utxoSet.FindUTXOPrefix(in.ID)
		if err != nil {
			return false
		}

		if len(outs.Outputs) == 0 || in.Out >= int64(len(outs.Outputs)) || in.Out < 0 {
			return false
		}
	}

	prevTxs := bc.GetTransaction(tx)

	return tx.Verify(prevTxs)
}

func (bc *Blockchain) MineBlock(transactions []*Transaction, address string, callback func([]*Transaction), ctx context.Context) (*Block, error) {
	var totalInput float64
	var totalOuput float64

	for _, tx := range transactions {
		publicKey := tx.Inputs[0].PubKey
		if !bc.VerifyTransaction(tx) {
			log.Error("Invalid Transaction")
			return nil, nil
		}

		for _, in := range tx.Inputs {
			if !bytes.Equal(publicKey, in.PubKey) {
				log.Error("Pubkey Does not math")
				return nil, nil
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

	nbits := bc.AdjustDifficulty(lastestBlock)

	reward, err := bc.GetBlockReward(int64(lastestBlock.Height+1), address)

	if err != nil {
		return nil, err
	}

	reward.Outputs[0].Value += totalInput - totalOuput

	transactions = append(transactions, reward)

	block, err := CreateBlock(transactions, lastestBlock.Hash, lastestBlock.Height+1, nbits, ctx)

	if err != nil {
		return nil, err
	}

	if block == nil {
		return nil, nil
	}

	err = bc.AddBlock(block, callback)

	if err != nil {
		return nil, fmt.Errorf("error adding block: %v", err)
	}

	return block, nil
}

func (bc *Blockchain) GetBestHeight() (int64, error) {
	var lastBlock Block

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(BestHeightPrefix))
		if err != nil {
			return err
		}

		lastHash, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		item, err = txn.Get(lastHash)
		if err != nil {
			return err
		}

		lastBlockData, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		block := DeserializeBlockData(lastBlockData)

		lastBlock = *block

		return nil

	})

	if err != nil {
		return 0, err
	}

	return lastBlock.Height, nil
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
