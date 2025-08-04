package blockchain

import (
	"bytes"
	"core-blockchain/common/utils"
	"encoding/hex"

	"github.com/dgraph-io/badger"
)

var (
	utxoPrefix = []byte("UTXO-")
	// prefixLength = len(utxoPrefix)
)

type UTXOSet struct {
	Blockchain *Blockchain
}

// func (u *UTXOSet) FindSpendableOutputs(publicKeyHash []byte, amount float64) (float64, map[string][]int) {
// 	unspentOuts := make(map[string][]int)
// 	accumulated := float64(0)

// 	db := u.Blockchain.Database

// 	err := db.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchValues = false

// 		it := txn.NewIterator(opts)
// 		defer it.Close()

// 		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
// 			item := it.Item()

// 			k := item.KeyCopy(nil)
// 			v, err := item.ValueCopy(nil)
// 			utils.ErrorHandle(err)

// 			outs := DeSerializeOuputs(v)

// 			k = bytes.TrimPrefix(k, utxoPrefix)

// 			txID := hex.EncodeToString(k)

// 			for outIdx, out := range outs.Outputs {
// 				if out.IsLockWithKey(publicKeyHash) && accumulated < amount {
// 					accumulated += out.Value
// 					unspentOuts[txID] = append(unspentOuts[txID], outIdx)

// 					if accumulated >= amount {
// 						break
// 					}
// 				}
// 			}

// 		}

// 		return nil
// 	})

// 	utils.ErrorHandle(err)

// 	return accumulated, unspentOuts
// }

func (u *UTXOSet) FindUTXOPrefix(txID []byte) (*TxOutputs, []byte) {
	var outs TxOutputs
	var key []byte

	err := u.Blockchain.Database.View(func(txn *badger.Txn) error {
		k := append(utxoPrefix, txID...)

		item, err := txn.Get([]byte(k))
		utils.ErrorHandle(err)

		key = item.KeyCopy(nil)

		v, err := item.ValueCopy(nil)
		utils.ErrorHandle(err)

		outs = DeSerializeOuputs(v)

		return nil

	})

	utils.ErrorHandle(err)

	return &outs, key
}

func (u *UTXOSet) FindSpendableOutputs(publicKeyHash []byte, amount float64) (float64, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := float64(0)

	iter := u.Blockchain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			outs, key := u.FindUTXOPrefix(tx.ID)

			k := bytes.TrimPrefix(key, utxoPrefix)

			txID := hex.EncodeToString(k)

			for outIdx, out := range outs.Outputs {
				if out.IsLockWithKey(publicKeyHash) && accumulated < amount {
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
					accumulated += out.Value

					if accumulated >= amount {
						break
					}
				}
			}

		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return accumulated, unspentOuts
}

// func (u *UTXOSet) FindUnSpentTransactions(pubKeyHash []byte) []TxOutput {
// 	var UTXOs []TxOutput
// 	db := u.Blockchain.Database

// 	err := db.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		it := txn.NewIterator(opts)
// 		defer it.Close()

// 		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
// 			item := it.Item()
// 			v, err := item.ValueCopy(nil)

// 			utils.ErrorHandle(err)
// 			outs := DeSerializeOuputs(v)

// 			for _, out := range outs.Outputs {
// 				if out.IsLockWithKey(pubKeyHash) {
// 					UTXOs = append(UTXOs, out)
// 				}
// 			}
// 		}

// 		return nil
// 	})

// 	utils.ErrorHandle(err)

// 	return UTXOs
// }

func (u *UTXOSet) FindUnSpentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput

	iter := u.Blockchain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			out, _ := u.FindUTXOPrefix(tx.ID)

			for _, out := range out.Outputs {
				if out.IsLockWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}

	}

	return UTXOs
}

func (u *UTXOSet) CountTransactions() int {
	db := u.Blockchain.Database
	counter := 0

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)

		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++
		}

		return nil
	})

	utils.ErrorHandle(err)

	return counter
}

func (u *UTXOSet) Update(bl *Block) {
	db := u.Blockchain.Database

	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range bl.Transactions {
			if tx.IsMinerTx() {
				newOutputs := TxOutputs{
					Outputs: tx.Outputs,
				}

				txID := append(utxoPrefix, tx.ID...)
				err := txn.Set(txID, newOutputs.Serialize())

				utils.ErrorHandle(err)
			} else {
				for _, in := range tx.Inputs {
					updatedOutputs := TxOutputs{}
					inID := append(utxoPrefix, tx.ID...)

					item, err := txn.Get(inID)
					utils.ErrorHandle(err)

					v, err := item.ValueCopy(nil)
					utils.ErrorHandle(err)

					outs := DeSerializeOuputs(v)

					for outIdx, out := range outs.Outputs {
						if int64(outIdx) != in.Out {
							updatedOutputs.Outputs = append(updatedOutputs.Outputs, out)
						}
					}

					if len(updatedOutputs.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							utils.ErrorHandle(err)
						}
					} else {
						if err := txn.Set(inID, updatedOutputs.Serialize()); err != nil {
							utils.ErrorHandle(err)
						}
					}
				}

				newOutputs := TxOutputs{
					Outputs: tx.Outputs,
				}

				txID := append(utxoPrefix, tx.ID...)

				err := txn.Set(txID, newOutputs.Serialize())

				utils.ErrorHandle(err)

			}
		}

		return nil
	})

	utils.ErrorHandle(err)
}

func (u *UTXOSet) Compute() {
	db := u.Blockchain.Database

	u.DeteleByPrefix(utxoPrefix)

	UTXO := u.Blockchain.FindUTXO()

	err := db.Update(func(txn *badger.Txn) error {
		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId)
			utils.ErrorHandle(err)

			key = append(utxoPrefix, key...)
			err = txn.Set(key, outs.Serialize())

			utils.ErrorHandle(err)
		}

		return nil
	})

	utils.ErrorHandle(err)

}

func (u *UTXOSet) DeteleByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := u.Blockchain.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}

		return nil
	}

	collectSize := 100000
	u.Blockchain.Database.View(func(txn *badger.Txn) error {
		otps := badger.DefaultIteratorOptions
		otps.PrefetchValues = false

		it := txn.NewIterator(otps)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++

			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					utils.ErrorHandle(err)
				}

				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				utils.ErrorHandle(err)
			}
		}

		return nil
	})
}
