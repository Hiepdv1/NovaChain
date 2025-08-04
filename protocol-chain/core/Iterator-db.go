package blockchain

import (
	"core-blockchain/common/utils"

	"github.com/dgraph-io/badger"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	if bc.LastHash == nil {
		return nil
	}

	return &BlockchainIterator{bc.LastHash, bc.Database}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utils.ErrorHandle(err)
		encodedBlock, err = item.ValueCopy(nil)
		utils.ErrorHandle(err)
		block = block.Deserialize(encodedBlock)

		return err
	})

	utils.ErrorHandle(err)

	iter.CurrentHash = block.PrevHash

	return block

}
