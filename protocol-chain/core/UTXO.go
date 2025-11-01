package blockchain

import (
	"bytes"
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

func (u *UTXOSet) FindUTXOPrefix(txID []byte) (*TxOutputs, []byte, error) {
	var outs TxOutputs
	var key []byte

	err := u.Blockchain.Database.View(func(txn *badger.Txn) error {
		k := append(utxoPrefix, txID...)

		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}

		key = item.KeyCopy(nil)

		v, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		outputs, err := DeSerializeOuputs(v)
		if err != nil {
			return err
		}

		outs = *outputs

		return nil

	})

	if err != nil {
		return nil, nil, err
	}

	return &outs, key, nil
}

func (u *UTXOSet) FindSpendableOutputs(publicKeyHash []byte, amount float64) (float64, map[string][]int, error) {
	unspentOuts := make(map[string][]int)
	accumulated := NewCoinAmountFromFloat(0.0)

	err := u.Blockchain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		opts.PrefetchValues = true
		it := txn.NewIterator(opts)

		defer it.Close()

		prefix := utxoPrefix

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			outs, err := DeSerializeOuputs(v)
			if err != nil {
				return err
			}

			key := item.Key()
			txID := hex.EncodeToString(bytes.TrimPrefix(key, prefix))

			for outIdx, out := range outs.Outputs {
				if out.IsLockWithKey(publicKeyHash) && accumulated.ToFloat() < amount {
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
					value := NewCoinAmountFromFloat(out.Value)
					accumulated = accumulated.Add(value)

					if accumulated.ToFloat() >= amount {
						break
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	return accumulated.ToFloat(), unspentOuts, nil
}

func (u *UTXOSet) FindUnSpentTransactions(pubKeyHash []byte) ([]TxOutput, error) {
	var UTXOs []TxOutput

	err := u.Blockchain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		opts.PrefetchValues = true
		it := txn.NewIterator(opts)

		defer it.Close()

		prefix := utxoPrefix

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			outs, err := DeSerializeOuputs(v)
			if err != nil {
				return err
			}

			for _, out := range outs.Outputs {
				if out.IsLockWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return UTXOs, nil
}

func (u *UTXOSet) CountTransactions() (int, error) {
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

	if err != nil {
		return 0, err
	}

	return counter, nil
}

func (u *UTXOSet) Update(bl *Block) error {
	db := u.Blockchain.Database

	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range bl.Transactions {
			if tx.IsMinerTx() {
				newOutputs := TxOutputs{
					Outputs: tx.Outputs,
				}

				txID := append(utxoPrefix, tx.ID...)
				serialize, err := newOutputs.Serialize()
				if err != nil {
					return err
				}

				err = txn.Set(txID, serialize)

				if err != nil {
					return err
				}
			} else {
				for _, in := range tx.Inputs {
					updatedOutputs := TxOutputs{}
					inID := append(utxoPrefix, tx.ID...)

					item, err := txn.Get(inID)
					if err != nil {
						return err
					}

					v, err := item.ValueCopy(nil)
					if err != nil {
						return err
					}

					outs, err := DeSerializeOuputs(v)
					if err != nil {
						return err
					}

					for outIdx, out := range outs.Outputs {
						if int64(outIdx) != in.Out {
							updatedOutputs.Outputs = append(updatedOutputs.Outputs, out)
						}
					}

					if len(updatedOutputs.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							return err
						}
					} else {
						serialize, err := updatedOutputs.Serialize()
						if err != nil {
							return err
						}

						if err := txn.Set(inID, serialize); err != nil {
							return err
						}
					}
				}

				newOutputs := TxOutputs{
					Outputs: tx.Outputs,
				}

				txID := append(utxoPrefix, tx.ID...)
				serialize, err := newOutputs.Serialize()
				if err != nil {
					return err
				}

				err = txn.Set(txID, serialize)

				if err != nil {
					return err
				}

			}
		}

		return nil
	})

	return err
}

func (u *UTXOSet) Compute() error {
	db := u.Blockchain.Database

	u.DeteleByPrefix(utxoPrefix)

	UTXO, _ := u.Blockchain.FindUTXO()

	err := db.Update(func(txn *badger.Txn) error {
		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId)
			if err != nil {
				return err
			}

			key = append(utxoPrefix, key...)
			serialize, err := outs.Serialize()
			if err != nil {
				return err
			}

			err = txn.Set(key, serialize)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err

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
					return err
				}

				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})
}
