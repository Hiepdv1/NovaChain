package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func StringPtrToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func (j *jobBlockSync) handleCreateInput(ins []transaction.TxInput, txHash string, tx *sql.Tx) error {
	ctx := context.Background()
	for _, in := range ins {
		pubkey, err := hex.DecodeString(in.PubKey)

		if err != nil {
			return err
		}

		args := dbchain.CreateTxInputParams{
			TxID:      txHash,
			InputTxID: in.ID,
			OutIndex:  in.Out,
			Sig:       in.Signature,
			PubKey:    in.PubKey,
		}
		txIn, err := j.dbTrans.CreateTxInput(ctx, args, tx)
		if err != nil {
			return err
		}

		log.Infof("Sync New TxInput: %v", txIn)

		if in.PubKey != "" {

			wallet, err := j.dbWallet.GetWalletByPubkey(ctx, []byte(in.PubKey), tx)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				log.Info("No wallet, Create Wallet in tx_input")
				newWallet, err := j.dbWallet.CreateWallet(
					ctx,
					dbwallet.CreateWalletParams{
						Address:       string(utils.PubKeyToAddress(pubkey)),
						PublicKey:     in.PubKey,
						PublicKeyHash: hex.EncodeToString(utils.PublicKeyHash(pubkey)),
						Balance:       "0",
						CreateAt: sql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						LastLogin: sql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
					},
					tx,
				)
				if err != nil {
					return err
				}
				wallet = &newWallet
			} else if err != nil {
				return err
			}

			txoutput, err := j.dbTrans.GetTxOutputByTxIDAndIndex(ctx, dbchain.GetTxOutputByTxIDAndIndexParams{
				TxID:  in.ID,
				Index: in.Out,
			})

			if err != nil && errors.Is(err, sql.ErrNoRows) {
				// block genesis or block reward
				txoutput = dbchain.TxOutput{
					Value: "0",
				}
			} else if err != nil {
				return err
			}

			value, err := strconv.ParseFloat(txoutput.Value, 64)
			if err != nil {
				return err
			}

			err = j.dbWallet.DecreaseWalletBalance(ctx, dbwallet.DecreaseWalletBalanceParams{
				Balance:   fmt.Sprintf("%.8f", value),
				Address:   wallet.Address,
				PublicKey: wallet.PublicKey,
			}, tx)

			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (j *jobBlockSync) handleCreateOutput(outs []transaction.TxOutput, txHash string, tx *sql.Tx) error {
	ctx := context.Background()
	for Index, out := range outs {
		wallet, err := j.dbWallet.GetWalletByPubKeyHash(ctx, out.PubKeyHash, tx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			log.Info("No wallet, Create Wallet in tx_outputs")
			newWallet, err := j.dbWallet.CreateWallet(
				ctx,
				dbwallet.CreateWalletParams{
					Address:       "-",
					PublicKey:     "-",
					PublicKeyHash: out.PubKeyHash,
					Balance:       "0",
					CreateAt: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					},
					LastLogin: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					},
				},
				tx,
			)
			if err != nil {
				return err
			}
			wallet = &newWallet
		} else if err != nil {
			return err
		}

		args := dbchain.CreateTxOutputParams{
			TxID:       txHash,
			Value:      fmt.Sprintf("%.8f", out.Value),
			PubKeyHash: out.PubKeyHash,
			Index:      int64(Index),
		}

		txout, err := j.dbTrans.CreateTxOutput(ctx, args, tx)
		if err != nil {
			return err
		}

		log.Infof("Sync new TxOutput: %v", txout)

		err = j.dbWallet.IncreaseWalletBalanceByPubKeyHash(
			ctx,
			dbwallet.IncreaseWalletBalanceByPubKeyHashParams{
				Balance:       fmt.Sprintf("%.8f", out.Value),
				PublicKeyHash: wallet.PublicKeyHash,
			},
			tx,
		)

		if err != nil {
			return err
		}

	}

	return nil
}

func (j *jobBlockSync) handleCreateTransactions(ctx context.Context, txs []*transaction.Transaction, hashBlock string, sqlTx *sql.Tx) error {
	for _, tx := range txs {
		args := dbchain.CreateTransactionParams{
			TxID: tx.ID,
			BID:  hashBlock,
		}

		transaction, err := j.dbTrans.CreateTransaction(ctx, args, sqlTx)
		if err != nil {
			return err
		}
		log.Infof("Sync New Transaction: %v", transaction)

		err = j.handleCreateInput(tx.Inputs, tx.ID, sqlTx)

		if err != nil {
			return err
		}

		err = j.handleCreateOutput(tx.Outputs, tx.ID, sqlTx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *jobBlockSync) handleCreateBlock(block *chain.Block, tx *sql.Tx) error {
	ctx := context.Background()

	args := dbchain.CreateBlockParams{
		BID:        block.Hash,
		PrevHash:   StringPtrToNullString(&block.PrevHash),
		Nonce:      block.Nonce,
		Height:     block.Height,
		MerkleRoot: block.MerkleRoot,
		Difficulty: block.Difficulty,
		TxCount:    block.TxCount,
		Timestamp:  block.Timestamp,
		NchainWork: block.NChainWork.String(),
	}

	newBlock, err := j.dbChain.CreateBlock(ctx, args, tx)

	if err != nil {
		return err
	}

	log.Infof("Block sync: %v", newBlock)

	return nil
}

func (j *jobBlockSync) isGenesisBlock(block *chain.Block) bool {
	return block.PrevHash == ""
}

func (j *jobBlockSync) handleReorganization(block dbchain.Block, sqlTx *sql.Tx) error {
	ctx := context.Background()
	lastBlock, err := j.dbChain.GetLastBlock(ctx)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil {
		return err
	}

	oldChain := make([]dbchain.Block, 0)

	currentBlock := lastBlock

	for {
		oldChain = append(oldChain, currentBlock)

		if currentBlock.BID == block.BID || !currentBlock.PrevHash.Valid {
			break
		}

		nextBlock, err := j.dbChain.GetBlockByHash(ctx, currentBlock.PrevHash.String)

		if err != nil {
			return err
		}

		currentBlock = nextBlock

	}

	for _, block := range oldChain {

		transactions, err := j.dbTrans.GetListTransactionByBlockHash(ctx, block.BID)
		if err != nil {
			return err
		}

		for _, tx := range transactions {

			txInputs, err := j.dbTrans.GetListTxInputByTxID(ctx, tx.TxID)
			if err != nil {
				return err
			}

			for _, input := range txInputs {
				if input.PubKey != "" {
					pubkey, err := hex.DecodeString(input.PubKey)
					if err != nil {
						return fmt.Errorf("failed to decode pubkey: %v", err)
					}

					wallet, err := j.dbWallet.GetWalletByPubkey(ctx, pubkey, sqlTx)
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						continue
					} else if err != nil {
						return err
					}

					prevOutput, err := j.dbTrans.GetTxOutputByTxIDAndIndex(ctx, dbchain.GetTxOutputByTxIDAndIndexParams{
						TxID:  input.InputTxID,
						Index: input.OutIndex,
					})

					if err != nil {
						return err
					}

					err = j.dbWallet.IncreaseWalletBalance(ctx, dbwallet.IncreaseWalletBalanceParams{
						Balance:   prevOutput.Value,
						Address:   wallet.Address,
						PublicKey: wallet.PublicKey,
					}, sqlTx)
					if err != nil {
						return err
					}
				}
			}

			txOutputs, err := j.dbTrans.GetListTxOutputByTxID(ctx, tx.TxID)
			if err != nil {
				return err
			}

			for _, output := range txOutputs {
				if output.PubKeyHash != "" {
					wallet, err := j.dbWallet.GetWalletByPubKeyHash(ctx, output.PubKeyHash, sqlTx)
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						continue
					} else if err != nil {
						return err
					}

					err = j.dbWallet.DecreaseWalletBalance(ctx, dbwallet.DecreaseWalletBalanceParams{
						Balance:   output.Value,
						Address:   wallet.Address,
						PublicKey: wallet.PublicKey,
					}, sqlTx)

					if err != nil {
						return err
					}
				}
			}
		}

		err = j.dbChain.DeleteBlockByHash(ctx, block.BID, sqlTx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *jobBlockSync) handleSyncBlock(blocks []*chain.Block, tx *sql.Tx) error {
	ctx := context.Background()

	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]

		lastBlock, err := j.dbChain.GetLastBlock(ctx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{}
		} else if err != nil {
			return err
		}

		if block.Hash == lastBlock.BID {
			log.Infof("Existing Block: %s", block.Hash)
			continue
		}

		if j.isGenesisBlock(block) {
			err := j.handleCreateBlock(block, tx)
			if err != nil {
				return err
			}
			err = j.handleCreateTransactions(ctx, block.Transactions, block.Hash, tx)

			if err != nil {
				return err
			}
			continue
		}

		prevBlock, err := j.dbChain.GetBlockByHash(ctx, block.PrevHash)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			log.Warnf("Block with hash %s not found, skipping", block.PrevHash)
			continue
		} else if err != nil {
			return err
		}

		lastBlockNChainWork := big.NewInt(0)
		lastBlockNChainWork.SetString(lastBlock.NchainWork, 10)

		if block.NChainWork.Cmp(lastBlockNChainWork) > 0 {
			if prevBlock.Height <= lastBlock.Height && prevBlock.BID != lastBlock.BID {
				err := j.handleReorganization(prevBlock, tx)
				if err != nil {
					return err
				}
			}

			err := j.handleCreateBlock(block, tx)
			if err != nil {
				return err
			}

			err = j.handleCreateTransactions(ctx, block.Transactions, block.Hash, tx)

			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (j *jobBlockSync) StartBlockSync(interval time.Duration) {
	log.Info("🔄 Chain sync started")
	ctx := context.Background()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {

		lastBlock, err := j.dbChain.GetLastBlock(ctx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{
				Height: 1,
			}
		} else if err != nil {
			log.Error("❌ GetLastBlock:", err)
			continue
		}

		log.Infof("⏱️ Sync tick: checking from height %d", lastBlock.Height)

		blocks, err := j.dbChainRpc.GetBlocksByHeightRange(lastBlock.Height, 10)
		if err != nil {
			log.Error("❌ Call RPC Error: ", err)
			continue
		}

		if len(blocks) == 0 {
			log.Warn("⚠️ No new blocks received")
			continue
		}

		log.Infof("📦 Received %d blocks", len(blocks))

		tx, err := db.Psql.BeginTx(ctx, nil)
		if err != nil {
			log.Error("❌ BeginTx error:", err)
			continue
		}

		err = j.handleSyncBlock(blocks, tx)
		if err != nil {
			log.Error("❌ Sync block error: ", err)
			if err := tx.Rollback(); err != nil {
				log.Error("❌ Rollback error: ", err)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Error("❌ Commit transaction error: ", err)
			continue
		}

		log.Infof("✅ Sync complete up to height: %d", lastBlock.Height)

	}

}
