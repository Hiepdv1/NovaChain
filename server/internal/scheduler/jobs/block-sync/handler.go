package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	dbutxo "ChainServer/internal/db/utxo"
	dbwallet "ChainServer/internal/db/wallet"
	"ChainServer/internal/states"
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

func (j *jobBlockSync) handleCreateUtxo(block *chain.Block, sqlTx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleCreateUtxo: Starting UTXO creation for block hash=%s height=%d txCount=%d", block.Hash, block.Height, block.TxCount)

	for _, tx := range block.Transactions {
		log.Debugf("handleCreateUtxo: Processing tx=%s in block=%s (isMiner=%v)", tx.ID, block.Hash, j.isTxMiner(tx))
		if !j.isTxMiner(tx) {
			for _, in := range tx.Inputs {
				log.Debugf("handleCreateUtxo: Deleting spent UTXO TxID=%s OutputIndex=%d for tx=%s in block=%s", in.ID, in.Out, tx.ID, block.Hash)
				params := dbutxo.DeleteUTXOParams{
					TxID:        in.ID,
					OutputIndex: in.Out,
				}

				if err := j.dbUtxo.DeleteUTXO(ctx, params, sqlTx); err != nil {
					log.Errorf("handleCreateUtxo: Failed to delete UTXO TxID=%s OutputIndex=%d for tx=%s in block=%s: %v", in.ID, in.Out, tx.ID, block.Hash, err)
					return err
				}
				log.Debugf("handleCreateUtxo: Successfully deleted UTXO TxID=%s OutputIndex=%d for tx=%s in block=%s", in.ID, in.Out, tx.ID, block.Hash)
			}
		}

		for idx, out := range tx.Outputs {
			log.Debugf("handleCreateUtxo: Creating UTXO for tx=%s output idx=%d value=%f pubKeyHash=%s in block=%s", tx.ID, idx, out.Value, out.PubKeyHash, block.Hash)
			params := dbutxo.CreateUTXOParams{
				TxID:        tx.ID,
				OutputIndex: int64(idx),
				Value:       fmt.Sprintf("%f", out.Value),
				PubKeyHash:  out.PubKeyHash,
				BlockID:     block.Hash,
			}

			if utxo, err := j.dbUtxo.CreateUTXO(ctx, params, sqlTx); err != nil {
				log.Errorf("handleCreateUtxo: Failed to create UTXO for tx=%s output idx=%d in block=%s: %v", tx.ID, idx, block.Hash, err)
				return err
			} else {
				log.Infof("handleCreateUtxo: New UTXO created TxID=%s OutputIndex=%d PubKeyHash=%s Value=%s in block=%s: %v", tx.ID, idx, out.PubKeyHash, utxo.Value, block.Hash, utxo)
			}
		}
	}

	log.Infof("handleCreateUtxo: Completed UTXO creation for block hash=%s (processed %d txs)", block.Hash, len(block.Transactions))
	return nil
}

func (j *jobBlockSync) handleCreateInput(ins []dto.TxInput, b_id, txHash string, tx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleCreateInput: Starting input creation for tx=%s block=%s (%d inputs)", txHash, b_id, len(ins))
	for _, in := range ins {
		log.Debugf("handleCreateInput: Processing input ID=%s Out=%d PubKey=%s Sig=%s for tx=%s block=%s", in.ID, in.Out, in.PubKey, in.Signature, txHash, b_id)
		pubkey, err := hex.DecodeString(in.PubKey)
		if err != nil {
			log.Errorf("handleCreateInput: Failed to decode pubkey=%s for input in tx=%s block=%s: %v", in.PubKey, txHash, b_id, err)
			return err
		}

		args := dbchain.CreateTxInputParams{
			TxID:      txHash,
			InputTxID: helpers.StringToNullString(in.ID),
			OutIndex:  in.Out,
			Sig:       helpers.StringToNullString(in.Signature),
			BID:       b_id,
			PubKey:    helpers.StringToNullString(in.PubKey),
		}
		txIn, err := j.dbTrans.CreateTxInput(ctx, args, tx)
		if err != nil {
			log.Errorf("handleCreateInput: Failed to create tx input for tx=%s block=%s: %v", txHash, b_id, err)
			return err
		}

		log.Infof("handleCreateInput: Created TxInput TxID=%s InputTxID=%s OutIndex=%d PubKey=%s Sig=%s BID=%s: %v", txHash, in.ID, in.Out, in.PubKey, in.Signature, b_id, txIn)

		if in.PubKey != "" {
			pubkeyHash := hex.EncodeToString(utils.PublicKeyHash(pubkey))
			wallet, err := j.dbWallet.GetWalletByPubKeyHash(ctx, pubkeyHash, tx)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				log.Infof("handleCreateInput: No wallet found for pubkey=%s in tx=%s block=%s, creating new wallet", in.PubKey, txHash, b_id)
				newWallet, err := j.dbWallet.CreateWallet(
					ctx,
					dbwallet.CreateWalletParams{
						Address:       helpers.StringToNullString(string(utils.PubKeyToAddress(pubkey))),
						PublicKey:     helpers.StringToNullString(in.PubKey),
						PublicKeyHash: pubkeyHash,
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
					log.Errorf("handleCreateInput: Failed to create wallet for pubkey=%s in tx=%s block=%s: %v", in.PubKey, txHash, b_id, err)
					return err
				}
				wallet = &newWallet
				log.Infof("handleCreateInput: Created new wallet Address=%s PublicKey=%s PublicKeyHash=%s for pubkey=%s in tx=%s block=%s", wallet.Address.String, wallet.PublicKey.String, wallet.PublicKeyHash, in.PubKey, txHash, b_id)
			} else if err != nil {
				log.Errorf("handleCreateInput: Failed to get wallet for pubkey=%s in tx=%s block=%s: %v", in.PubKey, txHash, b_id, err)
				return err
			}

			txoutput, err := j.dbTrans.GetTxOutputByTxIDAndIndex(ctx, dbchain.GetTxOutputByTxIDAndIndexParams{
				TxID:  in.ID,
				Index: in.Out,
			}, tx)

			if err != nil && errors.Is(err, sql.ErrNoRows) {
				// block genesis or block reward
				log.Warnf("handleCreateInput: No output found for input TxID=%s Index=%d in tx=%s block=%s (genesis/reward case)", in.ID, in.Out, txHash, b_id)
				txoutput = dbchain.TxOutput{
					Value: "0",
				}
			} else if err != nil {
				log.Errorf("handleCreateInput: Failed to get output for input TxID=%s Index=%d in tx=%s block=%s: %v", in.ID, in.Out, txHash, b_id, err)
				return err
			}

			log.Debugf("handleCreateInput: Decreasing wallet balance for wallet=%s by value=%s for input in tx=%s block=%s", wallet.Address.String, txoutput.Value, txHash, b_id)
			err = j.dbWallet.DecreaseWalletBalance(ctx, dbwallet.DecreaseWalletBalanceParams{
				Balance:   txoutput.Value,
				Address:   wallet.Address,
				PublicKey: wallet.PublicKey,
			}, tx)

			if err != nil {
				log.Errorf("handleCreateInput: Failed to decrease wallet balance for wallet=%s pubkey=%s in tx=%s block=%s by %s: %v", wallet.Address.String, in.PubKey, txHash, b_id, txoutput.Value, err)
				return err
			}
			log.Infof("handleCreateInput: Successfully decreased wallet balance for wallet=%s by %s in tx=%s block=%s", wallet.Address.String, txoutput.Value, txHash, b_id)
		}
	}

	log.Infof("handleCreateInput: Completed input creation for tx=%s block=%s", txHash, b_id)
	return nil
}

func (j *jobBlockSync) handleCreateOutput(outs []dto.TxOutput, b_id, txHash string, tx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleCreateOutput: Starting output creation for tx=%s block=%s (%d outputs)", txHash, b_id, len(outs))
	for Index, out := range outs {
		log.Debugf("handleCreateOutput: Processing output idx=%d value=%f pubKeyHash=%s for tx=%s block=%s", Index, out.Value, out.PubKeyHash, txHash, b_id)
		wallet, err := j.dbWallet.GetWalletByPubKeyHash(ctx, out.PubKeyHash, tx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			log.Infof("handleCreateOutput: No wallet found for PubKeyHash=%s in tx=%s output idx=%d block=%s, creating new wallet", out.PubKeyHash, txHash, Index, b_id)
			newWallet, err := j.dbWallet.CreateWallet(
				ctx,
				// Set null
				dbwallet.CreateWalletParams{
					Address:       helpers.StringToNullString(""),
					PublicKey:     helpers.StringToNullString(""),
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
				log.Errorf("handleCreateOutput: Failed to create wallet for PubKeyHash=%s in tx=%s output idx=%d block=%s: %v", out.PubKeyHash, txHash, Index, b_id, err)
				return err
			}
			wallet = &newWallet
			log.Infof("handleCreateOutput: Created new wallet PublicKeyHash=%s for output idx=%d in tx=%s block=%s", out.PubKeyHash, Index, txHash, b_id)
		} else if err != nil {
			log.Errorf("handleCreateOutput: Failed to get wallet for PubKeyHash=%s in tx=%s output idx=%d block=%s: %v", out.PubKeyHash, txHash, Index, b_id, err)
			return err
		}

		args := dbchain.CreateTxOutputParams{
			TxID:       txHash,
			Value:      fmt.Sprintf("%f", out.Value),
			PubKeyHash: out.PubKeyHash,
			Index:      int64(Index),
			BID:        b_id,
		}

		txout, err := j.dbTrans.CreateTxOutput(ctx, args, tx)
		if err != nil {
			log.Errorf("handleCreateOutput: Failed to create tx output for tx=%s index=%d block=%s: %v", txHash, Index, b_id, err)
			return err
		}

		log.Infof("handleCreateOutput: Created TxOutput TxID=%s Value=%s PubKeyHash=%s Index=%d BID=%s: %v", txHash, txout.Value, txout.PubKeyHash, txout.Index, b_id, txout)

		log.Debugf("handleCreateOutput: Increasing wallet balance for wallet=%s by value=%f for output idx=%d in tx=%s block=%s", wallet.Address.String, out.Value, Index, txHash, b_id)
		err = j.dbWallet.IncreaseWalletBalanceByPubKeyHash(
			ctx,
			dbwallet.IncreaseWalletBalanceByPubKeyHashParams{
				Balance:       fmt.Sprintf("%f", out.Value),
				PublicKeyHash: wallet.PublicKeyHash,
			},
			tx,
		)

		if err != nil {
			log.Errorf("handleCreateOutput: Failed to increase wallet balance for PubKeyHash=%s wallet=%s in tx=%s output idx=%d block=%s by %f: %v", out.PubKeyHash, wallet.Address.String, txHash, Index, b_id, out.Value, err)
			return err
		}
		log.Infof("handleCreateOutput: Successfully increased wallet balance for wallet=%s by %f in tx=%s output idx=%d block=%s", wallet.Address.String, out.Value, txHash, Index, b_id)
	}

	log.Infof("handleCreateOutput: Completed output creation for tx=%s block=%s", txHash, b_id)
	return nil
}

func (j *jobBlockSync) handleCreateTransactions(ctx context.Context, txs []*dto.Transaction, block *chain.Block, sqlTx *sql.Tx) error {
	log.Infof("handleCreateTransactions: Starting transaction creation for block hash=%s height=%d with %d txs", block.Hash, block.Height, len(txs))
	for _, tx := range txs {
		log.Debugf("handleCreateTransactions: Processing tx=%s in block=%s (isMiner=%v, inputs=%d, outputs=%d)", tx.ID, block.Hash, j.isTxMiner(tx), len(tx.Inputs), len(tx.Outputs))
		fromHash := ""
		toHash := ""
		amount := utils.NewCoinAmountFromFloat(0.0)
		fee := utils.NewCoinAmountFromFloat(0.0)

		if !j.isTxMiner(tx) {
			pubKeyBytes, err := hex.DecodeString(tx.Inputs[0].PubKey)
			if err != nil {
				log.Errorf("handleCreateTransactions: Failed to decode pubkey for first input in tx=%s block=%s: %v", tx.ID, block.Hash, err)
				return err
			}
			fromHash = hex.EncodeToString(utils.PublicKeyHash(pubKeyBytes))
			log.Debugf("handleCreateTransactions: Set fromHash=%s for tx=%s block=%s", fromHash, tx.ID, block.Hash)

			totalInput := utils.NewCoinAmountFromFloat(0.0)

			for _, in := range tx.Inputs {
				utxo, err := j.dbUtxo.GetUTXOByTxIDAndOut(ctx, dbutxo.GetUTXOByTxIDAndOutParams{
					TxID:        in.ID,
					OutputIndex: in.Out,
				}, sqlTx)

				if err != nil {
					log.Errorf("handleCreateTransactions: Failed to get UTXO for input TxID=%s Out=%d in tx=%s block=%s: %v", in.ID, in.Out, tx.ID, block.Hash, err)
					return err
				}

				value, err := strconv.ParseFloat(utxo.Value, 64)
				if err != nil {
					log.Errorf("handleCreateTransactions: Failed to parse UTXO value=%s for input TxID=%s Out=%d in tx=%s block=%s: %v", utxo.Value, in.ID, in.Out, tx.ID, block.Hash, err)
					return err
				}

				totalInput = totalInput.Add(utils.NewCoinAmountFromFloat(value))
				log.Debugf("handleCreateTransactions: Added input value=%f for TxID=%s Out=%d in tx=%s block=%s (totalInput=%f)", value, in.ID, in.Out, tx.ID, block.Hash, totalInput)
			}

			totalOutput := utils.NewCoinAmountFromFloat(0.0)

			for _, out := range tx.Outputs {
				if out.PubKeyHash != fromHash {
					toHash = out.PubKeyHash
					amount = amount.Add(utils.NewCoinAmountFromFloat(out.Value))
				}
				totalOutput = totalOutput.Add(utils.NewCoinAmountFromFloat(out.Value))
				log.Debugf("handleCreateTransactions: Processed output value=%f pubKeyHash=%s (toHash=%s amount=%f totalOutput=%f) in tx=%s block=%s", out.Value, out.PubKeyHash, toHash, amount, totalOutput, tx.ID, block.Hash)
			}

			fee = utils.SumFees(totalInput, totalOutput)
			log.Infof("handleCreateTransactions: Calculated for tx=%s block=%s: from=%s to=%s amount=%f fee=%f totalInput=%f totalOutput=%f", tx.ID, block.Hash, fromHash, toHash, amount, fee, totalInput, totalOutput)
		} else {
			log.Infof("handleCreateTransactions: Skipping fee/amount calc for miner tx=%s in block=%s", tx.ID, block.Hash)
		}

		args := dbchain.CreateTransactionParams{
			TxID:     tx.ID,
			BID:      block.Hash,
			Fromhash: helpers.StringToNullString(fromHash),
			Tohash:   helpers.StringToNullString(toHash),
			Amount:   helpers.FloatToNullString(amount.ToFloat()),
			Fee:      helpers.FloatToNullString(fee.ToFloat()),
			CreateAt: block.Timestamp,
		}

		transaction, err := j.dbTrans.CreateTransaction(ctx, args, sqlTx)
		if err != nil {
			log.Errorf("handleCreateTransactions: Failed to create transaction tx=%s for block=%s: %v", tx.ID, block.Hash, err)
			return err
		}
		log.Infof("handleCreateTransactions: Created Transaction tx=%s BID=%s From=%s To=%s Amount=%f Fee=%f CreateAt=%d: %v", tx.ID, block.Hash, fromHash, toHash, amount, fee, block.Timestamp, transaction)

		err = j.handleCreateInput(tx.Inputs, block.Hash, tx.ID, sqlTx)
		if err != nil {
			log.Errorf("handleCreateTransactions: Failed to create inputs for tx=%s block=%s: %v", tx.ID, block.Hash, err)
			return err
		}

		err = j.handleCreateOutput(tx.Outputs, block.Hash, tx.ID, sqlTx)
		if err != nil {
			log.Errorf("handleCreateTransactions: Failed to create outputs for tx=%s block=%s: %v", tx.ID, block.Hash, err)
			return err
		}

		exists, err := j.dbTrans.ExistingPendingTransaction(
			ctx,
			dbPendingTx.PendingTxExistsParams{
				TxID: tx.ID,
				Status: []string{
					string(constants.TxStatusPending),
				},
			},
			sqlTx,
		)

		if err != nil {
			log.Errorf("handleCreateTransactions: Failed to check pending tx=%s in block=%s: %v", tx.ID, block.Hash, err)
			return err
		}

		if exists {
			log.Infof("handleCreateTransactions: Updating pending tx=%s to mined status for block=%s", tx.ID, block.Hash)
			_, err := j.dbTrans.UpdatePendingTxsStatus(ctx, dbPendingTx.UpdatePendingTxsStatusParams{
				NewStatus: string(constants.TxStatusMined),
				TxIds:     []string{tx.ID},
				OldStatus: []string{
					string(constants.TxStatusPending),
				},
			}, sqlTx)

			if err != nil {
				log.Errorf("handleCreateTransactions: Failed to update pending tx=%s to mined for block=%s: %v", tx.ID, block.Hash, err)
				return err
			}
			log.Infof("handleCreateTransactions: Successfully updated pending tx=%s to mined for block=%s", tx.ID, block.Hash)
		} else {
			log.Debugf("handleCreateTransactions: No pending tx found for tx=%s in block=%s", tx.ID, block.Hash)
		}
	}

	log.Infof("handleCreateTransactions: Completed transaction creation for block hash=%s (processed %d txs)", block.Hash, len(txs))
	return nil
}

func (j *jobBlockSync) handleCreateBlock(block *chain.Block, tx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleCreateBlock: Starting block creation for hash=%s height=%d prevHash=%s txCount=%d", block.Hash, block.Height, block.PrevHash, block.TxCount)

	size, err := utils.GobEncode(block)
	if err != nil {
		log.Errorf("handleCreateBlock: Failed to gob encode block=%s: %v", block.Hash, err)
		return err
	}

	args := dbchain.CreateBlockParams{
		BID:        block.Hash,
		PrevHash:   helpers.StringPtrToNullString(&block.PrevHash),
		Nonce:      block.Nonce,
		Height:     block.Height,
		MerkleRoot: block.MerkleRoot,
		Nbits:      block.NBits,
		TxCount:    block.TxCount,
		Timestamp:  block.Timestamp,
		NchainWork: block.NChainWork.String(),
		Size:       float64(len(size)),
	}

	newBlock, err := j.dbChain.CreateBlock(ctx, args, tx)

	if err != nil {
		log.Errorf("handleCreateBlock: Failed to create block=%s: %v", block.Hash, err)
		return err
	}

	log.Infof("handleCreateBlock: Successfully created block BID=%s Height=%d PrevHash=%s Nonce=%d MerkleRoot=%s Nbits=%d TxCount=%d Timestamp=%d NchainWork=%s Size=%.2fMB: %v", newBlock.BID, newBlock.Height, newBlock.PrevHash.String, newBlock.Nonce, newBlock.MerkleRoot, newBlock.Nbits, newBlock.TxCount, newBlock.Timestamp, newBlock.NchainWork, newBlock.Size, newBlock)
	return nil
}

func (j *jobBlockSync) isGenesisBlock(block *chain.Block) bool {
	return block.PrevHash == ""
}

func (j *jobBlockSync) isTxMiner(tx *dto.Transaction) bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}
func (j *jobBlockSync) handleReorganizationBlocks(tx *sql.Tx) error {
	log.Info("🔁 [Reorg] Starting chain reorganization...")
	ctx := context.Background()

	locator, err := j.dbChain.GetBlockLocator(ctx, tx)
	if err != nil {
		log.Error("❌ [Reorg] Failed to get block locator: ", err)
		return err
	}

	bestHeight, err := j.dbChain.GetBestHeight(ctx, tx)
	if err != nil {
		log.Error("❌ [Reorg] Failed to get best height: ", err)
		return err
	}

	commonBlock, err := j.dbChainRpc.GetCommonBlock(locator)
	if err != nil {
		log.Error("❌ [Reorg] Failed to get common block from RPC: ", err)
		return err
	}

	log.Infof("🔍 [Reorg] Common block found at height=%d (hash=%s)", commonBlock.Height, commonBlock.Hash)
	log.Infof("🧹 [Reorg] Removing blocks from height=%d to height=%d...", commonBlock.Height, bestHeight)

	err = j.dbChain.DeleteBlockByRangeHeight(ctx, commonBlock.Height+1, bestHeight, tx)
	if err != nil {
		log.Error("❌ [Reorg] Failed to delete old blocks: ", err)
		return err
	}

	log.Info("✅ [Reorg] Chain reorganization completed successfully.")
	return nil
}

func (j *jobBlockSync) handleSyncBlock(blocks []*chain.Block, tx *sql.Tx) error {
	ctx := context.Background()
	start := blocks[0].Height
	end := blocks[len(blocks)-1].Height
	log.Infof("📦 [Sync] Starting block sync (%d blocks | heights %d → %d)", len(blocks), start, end)

	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]
		log.Infof("➡️ [Sync] Processing block #%d (height=%d | hash=%s)", len(blocks)-i, block.Height, block.Hash)

		lastBlock, err := j.dbChain.GetLastBlock(ctx, tx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{}
			log.Debugf("ℹ️ [Sync] No last block found — treating as empty chain (block=%s height=%d)", block.Hash, block.Height)
		} else if err != nil {
			log.Errorf("❌ [Sync] Failed to get last block (block=%s height=%d): %v", block.Hash, block.Height, err)
			return err
		}

		exists, err := j.dbChain.ExistingBlock(ctx, block.Hash, tx)
		if err != nil {
			log.Errorf("❌ [Sync] Failed to check block existence (hash=%s): %v", block.Hash, err)
			return err
		}
		if exists {
			log.Debugf("⚠️ [Sync] Block already exists (hash=%s height=%d), skipping.", block.Hash, block.Height)
			continue
		}

		if j.isGenesisBlock(block) {
			log.Infof("🌱 [Sync] Found genesis block (hash=%s)", block.Hash)

			if err := j.handleCreateBlock(block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to create genesis block (hash=%s): %v", block.Hash, err)
				return err
			}
			if err := j.handleCreateTransactions(ctx, block.Transactions, block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to create genesis transactions (hash=%s): %v", block.Hash, err)
				return err
			}
			if err := j.handleCreateUtxo(block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to create genesis UTXO (hash=%s): %v", block.Hash, err)
				return err
			}

			log.Info("✅ [Sync] Genesis block processed successfully.")
			continue
		}

		lastBlockWork := big.NewInt(0)
		lastBlockWork.SetString(lastBlock.NchainWork, 10)

		if block.NChainWork.Cmp(lastBlockWork) > 0 {
			log.Infof("🔗 [Sync] Stronger block found (hash=%s height=%d)", block.Hash, block.Height)

			if _, err := j.dbChain.GetBlockByHash(ctx, block.PrevHash, tx); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Warnf("⚠️ [Sync] Missing previous block (prevHash=%s) → triggering reorganization", block.PrevHash)
					if err := j.handleReorganizationBlocks(tx); err != nil {
						log.Errorf("❌ [Sync] Reorganization failed (block=%s height=%d): %v", block.Hash, block.Height, err)
						return err
					}
					break
				}
				log.Errorf("❌ [Sync] Failed to fetch previous block (prevHash=%s): %v", block.PrevHash, err)
				return err
			}

			if err := j.handleCreateBlock(block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to store block (hash=%s): %v", block.Hash, err)
				return err
			}
			if err := j.handleCreateTransactions(ctx, block.Transactions, block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to store block transactions (hash=%s): %v", block.Hash, err)
				return err
			}
			if err := j.handleCreateUtxo(block, tx); err != nil {
				log.Errorf("❌ [Sync] Failed to store UTXOs (hash=%s): %v", block.Hash, err)
				return err
			}

			log.Infof("✅ [Sync] Block added successfully (height=%d hash=%s)", block.Height, block.Hash)
		} else {
			log.Warnf("⚠️ [Sync] Weaker chain detected (hash=%s height=%d) — skipping.", block.Hash, block.Height)
		}
	}

	log.Infof("🏁 [Sync] Completed block sync (%d blocks processed)", len(blocks))
	return nil
}

func (j *jobBlockSync) StartBlockSync(interval time.Duration) {
	log.Info("🚀 [SyncLoop] Chain sync started")
	ctx := context.Background()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		lastBlock, err := j.dbChain.GetLastBlock(ctx, nil)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{Height: 1}
			log.Warn("⚠️ [SyncLoop] No last block found — starting from height=1")
		} else if err != nil {
			log.Errorf("❌ [SyncLoop] Failed to get last block: %v", err)
			continue
		}

		log.Infof("🕒 [SyncLoop] Tick: currentHeight=%d → fetching next 100 blocks...", lastBlock.Height)

		blocks, err := j.dbChainRpc.GetBlocksByHeightRange(lastBlock.Height, env.Cfg.Sync_block_batch_size)
		if err != nil {
			log.Errorf("❌ [SyncLoop] RPC error (startHeight=%d): %v", lastBlock.Height, err)
			continue
		}
		if len(blocks) == 0 {
			log.Info("✅ [SyncLoop] No new blocks — chain is up-to-date.")
			continue
		}

		log.Infof("📥 [SyncLoop] Retrieved %d new blocks (heights %d → %d)", len(blocks), blocks[len(blocks)-1].Height, blocks[0].Height)

		tx, err := db.Psql.BeginTx(ctx, nil)
		if err != nil {
			log.Errorf("❌ [SyncLoop] Failed to start DB transaction: %v", err)
			continue
		}

		log.Debugf("🧩 [SyncLoop] Transaction started for %d blocks", len(blocks))

		if err := j.handleSyncBlock(blocks, tx); err != nil {
			log.Errorf("❌ [SyncLoop] Sync failed: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Errorf("❌ [SyncLoop] Rollback failed: %v", rbErr)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Errorf("❌ [SyncLoop] Commit failed: %v", err)
			continue
		}

		bestHeight, err := j.dbChain.GetBestHeight(ctx, nil)
		if err != nil {
			log.Errorf("⚠️ [SyncLoop] Failed to fetch best height: %v", err)
		} else {
			states.ChainSyncState.SyncStatus = utils.EvaluateSyncStatus(bestHeight, blocks[0].Height)
		}

		log.Infof("🎯 [SyncLoop] Sync successful → bestHeight=%d (processed %d blocks)", bestHeight, len(blocks))
	}
}
