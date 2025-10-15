package blockchain

import (
	"github.com/dgraph-io/badger"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *Blockchain) Iterator() (*BlockchainIterator, error) {
	LastHash, err := bc.GetLastBlock()
	if err != nil {
		return nil, err
	}

	return &BlockchainIterator{LastHash.Hash, bc.Database}, nil
}

func (iter *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return err
		}
		encodedBlock, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		block = DeserializeBlockData(encodedBlock)

		return err
	})

	if err != nil {
		return nil, err
	}

	iter.CurrentHash = block.PrevHash

	return block, nil

}
