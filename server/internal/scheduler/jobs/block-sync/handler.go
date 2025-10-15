package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	dbutxo "ChainServer/internal/db/utxo"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
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
				Value:       helpers.FormatDecimal(out.Value, 10),
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

func (j *jobBlockSync) handleReorganizationUtxo(block *dbchain.Block, sqlTx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleReorganizationUtxo: Starting UTXO reorganization for block hash=%s height=%d", block.BID, block.Height)

	txOutputs, err := j.dbTrans.FindListTxOutputByBlockHash(ctx, block.BID, sqlTx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			txOutputs = make([]dbchain.TxOutput, 0)
			log.Debugf("handleReorganizationUtxo: No tx outputs found for block=%s", block.BID)
		} else {
			log.Errorf("handleReorganizationUtxo: Failed to fetch tx outputs for block=%s: %v", block.BID, err)
			return err
		}
	}

	txInputs, err := j.dbTrans.FindListTxInputByBlockHash(ctx, block.BID, sqlTx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			txInputs = make([]dbchain.TxInput, 0)
			log.Debugf("handleReorganizationUtxo: No tx inputs found for block=%s", block.BID)
		} else {
			log.Errorf("handleReorganizationUtxo: Failed to fetch tx inputs for block=%s: %v", block.BID, err)
			return err
		}
	}

	log.Infof("handleReorganizationUtxo: Deleting %d UTXOs for block=%s during reorganization", len(txOutputs), block.BID)
	for range txOutputs {
		if err := j.dbUtxo.DeleteUTXOByBlockID(ctx, block.BID, sqlTx); err != nil {
			log.Errorf("handleReorganizationUtxo: Failed to delete UTXO for block=%s: %v", block.BID, err)
			return err
		}
		log.Debugf("handleReorganizationUtxo: Deleted UTXO for block=%s", block.BID)
	}

	for _, in := range txInputs {
		log.Debugf("handleReorganizationUtxo: Recreating UTXO for input TxID=%s OutIndex=%d in block=%s", in.InputTxID.String, in.OutIndex, block.BID)
		getOutparams := dbchain.GetTxOutputByTxIDAndIndexParams{
			TxID:  in.InputTxID.String,
			Index: in.OutIndex,
		}
		output, err := j.dbTrans.GetTxOutputByTxIDAndIndex(ctx, getOutparams, sqlTx)
		if err != nil {
			log.Errorf("handleReorganizationUtxo: Failed to get output for input TxID=%s Index=%d in block=%s: %v", in.InputTxID.String, in.OutIndex, block.BID, err)
			return err
		}

		createUtxoParams := dbutxo.CreateUTXOParams{
			TxID:        in.InputTxID.String,
			OutputIndex: in.OutIndex,
			Value:       output.Value,
			PubKeyHash:  output.PubKeyHash,
			BlockID:     block.BID,
		}

		utxo, err := j.dbUtxo.CreateUTXO(ctx, createUtxoParams, sqlTx)
		if err != nil {
			log.Errorf("handleReorganizationUtxo: Failed to recreate UTXO for input TxID=%s Index=%d in block=%s: %v", in.InputTxID.String, in.OutIndex, block.BID, err)
			return err
		}

		log.Infof("handleReorganizationUtxo: Recreated UTXO TxID=%s OutputIndex=%d Value=%s PubKeyHash=%s in block=%s: %v", in.InputTxID.String, in.OutIndex, output.Value, output.PubKeyHash, block.BID, utxo)
	}

	log.Infof("handleReorganizationUtxo: Completed UTXO reorganization for block hash=%s (recreated %d UTXOs)", block.BID, len(txInputs))
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

			value, err := strconv.ParseFloat(txoutput.Value, 64)
			if err != nil {
				log.Errorf("handleCreateInput: Failed to parse value=%s for input in tx=%s block=%s: %v", txoutput.Value, txHash, b_id, err)
				return err
			}

			log.Debugf("handleCreateInput: Decreasing wallet balance for wallet=%s by value=%f for input in tx=%s block=%s", wallet.Address.String, value, txHash, b_id)
			err = j.dbWallet.DecreaseWalletBalance(ctx, dbwallet.DecreaseWalletBalanceParams{
				Balance:   fmt.Sprintf("%.8f", value),
				Address:   wallet.Address,
				PublicKey: wallet.PublicKey,
			}, tx)

			if err != nil {
				log.Errorf("handleCreateInput: Failed to decrease wallet balance for wallet=%s pubkey=%s in tx=%s block=%s by %f: %v", wallet.Address.String, in.PubKey, txHash, b_id, value, err)
				return err
			}
			log.Infof("handleCreateInput: Successfully decreased wallet balance for wallet=%s by %f in tx=%s block=%s", wallet.Address.String, value, txHash, b_id)
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
			Value:      fmt.Sprintf("%.8f", out.Value),
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
				Balance:       fmt.Sprintf("%.8f", out.Value),
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
		amount := 0.0
		fee := 0.0

		if !j.isTxMiner(tx) {
			pubKeyBytes, err := hex.DecodeString(tx.Inputs[0].PubKey)
			if err != nil {
				log.Errorf("handleCreateTransactions: Failed to decode pubkey for first input in tx=%s block=%s: %v", tx.ID, block.Hash, err)
				return err
			}
			fromHash = hex.EncodeToString(utils.PublicKeyHash(pubKeyBytes))
			log.Debugf("handleCreateTransactions: Set fromHash=%s for tx=%s block=%s", fromHash, tx.ID, block.Hash)

			totalInput := 0.0

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

				totalInput += value
				log.Debugf("handleCreateTransactions: Added input value=%f for TxID=%s Out=%d in tx=%s block=%s (totalInput=%f)", value, in.ID, in.Out, tx.ID, block.Hash, totalInput)
			}

			totalOutput := 0.0

			for _, out := range tx.Outputs {
				if out.PubKeyHash != fromHash {
					toHash = out.PubKeyHash
					amount += out.Value
				}
				totalOutput += out.Value
				log.Debugf("handleCreateTransactions: Processed output value=%f pubKeyHash=%s (toHash=%s amount=%f totalOutput=%f) in tx=%s block=%s", out.Value, out.PubKeyHash, toHash, amount, totalOutput, tx.ID, block.Hash)
			}

			fee = totalInput - totalOutput
			log.Infof("handleCreateTransactions: Calculated for tx=%s block=%s: from=%s to=%s amount=%f fee=%f totalInput=%f totalOutput=%f", tx.ID, block.Hash, fromHash, toHash, amount, fee, totalInput, totalOutput)
		} else {
			log.Infof("handleCreateTransactions: Skipping fee/amount calc for miner tx=%s in block=%s", tx.ID, block.Hash)
		}

		args := dbchain.CreateTransactionParams{
			TxID:     tx.ID,
			BID:      block.Hash,
			Fromhash: helpers.StringToNullString(fromHash),
			Tohash:   helpers.StringToNullString(toHash),
			Amount:   helpers.FloatToNullString(amount),
			Fee:      helpers.FloatToNullString(fee),
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
					string(constants.TxStatusMining),
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
					string(constants.TxStatusMining),
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

func (j *jobBlockSync) handleReorganization(block *chain.Block, sqlTx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleReorganization: Starting reorganization triggered by new block hash=%s height=%d", block.Hash, block.Height)
	lastBlock, err := j.dbChain.GetLastBlock(ctx, sqlTx)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Warnf("handleReorganization: No last block found during reorganization for new block=%s", block.Hash)
		return err
	} else if err != nil {
		log.Errorf("handleReorganization: Failed to get last block during reorganization for new block=%s: %v", block.Hash, err)
		return err
	}
	log.Debugf("handleReorganization: Current last block hash=%s height=%d", lastBlock.BID, lastBlock.Height)

	oldChain := make([]dbchain.Block, 0)

	currentBlock := lastBlock

	for {
		oldChain = append(oldChain, currentBlock)
		log.Debugf("handleReorganization: Added to old chain: block=%s height=%d prevHash=%s", currentBlock.BID, currentBlock.Height, currentBlock.PrevHash.String)

		nextBlock, err := j.dbChain.GetBlockByHash(ctx, currentBlock.PrevHash.String, sqlTx)

		if err != nil {
			log.Errorf("handleReorganization: Failed to get prev block=%s during reorganization for new block=%s: %v", currentBlock.PrevHash.String, block.Hash, err)
			return err
		}

		if nextBlock.BID == block.PrevHash || !currentBlock.PrevHash.Valid {
			log.Debugf("handleReorganization: Reached fork point at block=%s (new block prev=%s)", nextBlock.BID, block.PrevHash)
			break
		}

		currentBlock = nextBlock
	}

	log.Infof("handleReorganization: Identified %d blocks to rollback for reorganization (new block=%s): %v", len(oldChain), block.Hash, oldChain)

	for i, blk := range oldChain {
		log.Infof("handleReorganization: Rolling back block %d/%d: hash=%s height=%d in reorganization (new block=%s)", i+1, len(oldChain), blk.BID, blk.Height, block.Hash)

		if err := j.handleReorganizationUtxo(&blk, sqlTx); err != nil {
			log.Errorf("handleReorganization: Failed to reorganize UTXOs for rollback block=%s height=%d (new block=%s): %v", blk.BID, blk.Height, block.Hash, err)
			return err
		}
		log.Infof("handleReorganization: UTXOs reorganized for rollback block=%s height=%d", blk.BID, blk.Height)

		transactions, err := j.dbTrans.GetListTransactionByBlockHash(ctx, blk.BID, sqlTx)
		if err != nil {
			log.Errorf("handleReorganization: Failed to get %d transactions for rollback block=%s height=%d (new block=%s): %v", len(transactions), blk.BID, blk.Height, block.Hash, err)
			return err
		}
		log.Debugf("handleReorganization: Found %d txs to rollback in block=%s height=%d", len(transactions), blk.BID, blk.Height)

		for _, tx := range transactions {
			log.Infof("handleReorganization: Rolling back tx=%s in rollback block=%s height=%d (new block=%s)", tx.TxID, blk.BID, blk.Height, block.Hash)

			txInputs, err := j.dbTrans.GetListTxInputByTxID(ctx, tx.TxID)
			if err != nil {
				log.Errorf("handleReorganization: Failed to get inputs for tx=%s in rollback block=%s height=%d (new block=%s): %v", tx.TxID, blk.BID, blk.Height, block.Hash, err)
				return err
			}

			for _, input := range txInputs {
				if input.PubKey.Valid {
					pubkey, err := hex.DecodeString(input.PubKey.String)
					if err != nil {
						log.Errorf("handleReorganization: Failed to decode pubkey=%s for input in tx=%s rollback block=%s height=%d (new block=%s): %v", input.PubKey.String, tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return fmt.Errorf("failed to decode pubkey: %v", err)
					}

					wallet, err := j.dbWallet.GetWalletByPubkey(ctx, pubkey, sqlTx)
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						log.Warnf("handleReorganization: No wallet found for pubkey in input tx=%s rollback block=%s height=%d (new block=%s), skipping balance restore", tx.TxID, blk.BID, blk.Height, block.Hash)
						continue
					} else if err != nil {
						log.Errorf("handleReorganization: Failed to get wallet for pubkey in input tx=%s rollback block=%s height=%d (new block=%s): %v", tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return err
					}

					prevOutput, err := j.dbTrans.GetTxOutputByTxIDAndIndex(ctx, dbchain.GetTxOutputByTxIDAndIndexParams{
						TxID:  input.InputTxID.String,
						Index: input.OutIndex,
					}, sqlTx)

					if err != nil {
						log.Errorf("handleReorganization: Failed to get prev output for input TxID=%s Index=%d in tx=%s rollback block=%s height=%d (new block=%s): %v", input.InputTxID.String, input.OutIndex, tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return err
					}

					log.Debugf("handleReorganization: Restoring balance %s for wallet=%s in input tx=%s rollback block=%s height=%d", prevOutput.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height)
					err = j.dbWallet.IncreaseWalletBalance(ctx, dbwallet.IncreaseWalletBalanceParams{
						Balance:   prevOutput.Value,
						Address:   wallet.Address,
						PublicKey: wallet.PublicKey,
					}, sqlTx)
					if err != nil {
						log.Errorf("handleReorganization: Failed to increase wallet balance %s for wallet=%s in input tx=%s rollback block=%s height=%d (new block=%s): %v", prevOutput.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return err
					}
					log.Infof("handleReorganization: Restored balance %s for wallet=%s in input tx=%s rollback block=%s height=%d", prevOutput.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height)
				}
			}

			txOutputs, err := j.dbTrans.GetListTxOutputByTxID(ctx, tx.TxID)
			if err != nil {
				log.Errorf("handleReorganization: Failed to get outputs for tx=%s in rollback block=%s height=%d (new block=%s): %v", tx.TxID, blk.BID, blk.Height, block.Hash, err)
				return err
			}

			for _, output := range txOutputs {
				if output.PubKeyHash != "" {
					wallet, err := j.dbWallet.GetWalletByPubKeyHash(ctx, output.PubKeyHash, sqlTx)
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						log.Warnf("handleReorganization: No wallet found for PubKeyHash=%s in output tx=%s rollback block=%s height=%d (new block=%s), skipping balance deduct", output.PubKeyHash, tx.TxID, blk.BID, blk.Height, block.Hash)
						continue
					} else if err != nil {
						log.Errorf("handleReorganization: Failed to get wallet for PubKeyHash=%s in output tx=%s rollback block=%s height=%d (new block=%s): %v", output.PubKeyHash, tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return err
					}

					log.Debugf("handleReorganization: Deducting balance %s for wallet=%s in output tx=%s rollback block=%s height=%d", output.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height)
					err = j.dbWallet.DecreaseWalletBalance(ctx, dbwallet.DecreaseWalletBalanceParams{
						Balance:   output.Value,
						Address:   wallet.Address,
						PublicKey: wallet.PublicKey,
					}, sqlTx)

					if err != nil {
						log.Errorf("handleReorganization: Failed to decrease wallet balance %s for wallet=%s in output tx=%s rollback block=%s height=%d (new block=%s): %v", output.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height, block.Hash, err)
						return err
					}
					log.Infof("handleReorganization: Deducted balance %s for wallet=%s in output tx=%s rollback block=%s height=%d", output.Value, wallet.Address.String, tx.TxID, blk.BID, blk.Height)
				}
			}
		}

		err = j.dbChain.DeleteBlockByHash(ctx, blk.BID, sqlTx)
		if err != nil {
			log.Errorf("handleReorganization: Failed to delete rollback block=%s height=%d (new block=%s): %v", blk.BID, blk.Height, block.Hash, err)
			return err
		}
		log.Infof("handleReorganization: Deleted rollback block=%s height=%d", blk.BID, blk.Height)
	}

	log.Infof("handleReorganization: Completed reorganization for new block hash=%s (rolled back %d blocks)", block.Hash, len(oldChain))
	return nil
}

func (j *jobBlockSync) handleSwitchChain(block *chain.Block, tx *sql.Tx) error {
	log.Infof("handleSwitchChain: Starting chain switch triggered by block hash=%s height=%d nChainWork=%s", block.Hash, block.Height, block.NChainWork)
	newChain := make([]*chain.Block, 0)
	newChain = append(newChain, block)
	currentHash := block.PrevHash

	for {
		blk, err := j.dbChainRpc.GetBlockByHash(currentHash)
		if err != nil {
			log.Errorf("handleSwitchChain: Failed to fetch block by hash=%s during chain switch (new block=%s): %v", currentHash, block.Hash, err)
			return err
		}

		exists, err := j.dbChain.ExistingBlock(context.Background(), blk.Hash, nil)
		if err != nil {
			log.Errorf("handleSwitchChain: Failed to check existence of block hash=%s during chain switch (new block=%s): %v", blk.Hash, block.Hash, err)
			return err
		}

		if exists {
			log.Infof("handleSwitchChain: Existing block hash=%s found during chain switch, stopping fetch (new block=%s)", blk.Hash, block.Hash)
			break
		}

		newChain = append(newChain, blk)
		currentHash = blk.PrevHash
		log.Infof("handleSwitchChain: Added block hash=%s height=%d to new chain during switch (new block=%s)", blk.Hash, blk.Height, block.Hash)
	}

	slices.Reverse(newChain)
	log.Infof("handleSwitchChain: New chain prepared with %d blocks for switch (first=%s last=%s new block=%s)", len(newChain), newChain[0].Hash, newChain[len(newChain)-1].Hash, block.Hash)

	err := j.handleReorganization(newChain[0], tx)
	if err != nil {
		log.Errorf("handleSwitchChain: Failed to reorganize for chain switch (new block=%s): %v", block.Hash, err)
		return err
	}
	log.Infof("handleSwitchChain: Reorganization completed for chain switch (new block=%s)", block.Hash)

	log.Infof("handleSwitchChain: Switching to new chain starting with hash=%s nChainWork=%s", block.Hash, block.NChainWork)

	for i, blk := range newChain {
		log.Infof("handleSwitchChain: Creating block %d/%d in new chain: hash=%s height=%d (new block=%s)", i+1, len(newChain), blk.Hash, blk.Height, block.Hash)
		err := j.handleCreateBlock(blk, tx)
		if err != nil {
			log.Errorf("handleSwitchChain: Failed to create block hash=%s height=%d in new chain (new block=%s): %v", blk.Hash, blk.Height, block.Hash, err)
			return err
		}
		err = j.handleCreateUtxo(blk, tx)
		if err != nil {
			log.Errorf("handleSwitchChain: Failed to create UTXOs for block hash=%s height=%d in new chain (new block=%s): %v", blk.Hash, blk.Height, block.Hash, err)
			return err
		}
		log.Infof("handleSwitchChain: Successfully created block hash=%s height=%d and UTXOs in new chain", blk.Hash, blk.Height)
	}

	log.Infof("handleSwitchChain: Completed chain switch for block hash=%s (added %d new blocks)", block.Hash, len(newChain))
	return nil
}

func (j *jobBlockSync) handleSyncBlock(blocks []*chain.Block, tx *sql.Tx) error {
	ctx := context.Background()
	log.Infof("handleSyncBlock: Starting sync for %d blocks (range heights %d-%d)", len(blocks), blocks[0].Height, blocks[len(blocks)-1].Height)

	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]
		log.Infof("handleSyncBlock: Processing block %d/%d: hash=%s height=%d prevHash=%s nChainWork=%s", i+1, len(blocks), block.Hash, block.Height, block.PrevHash, block.NChainWork)

		lastBlock, err := j.dbChain.GetLastBlock(ctx, tx)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{}
			log.Debugf("handleSyncBlock: No last block in DB, treating as empty chain for block=%s height=%d", block.Hash, block.Height)
		} else if err != nil {
			log.Errorf("handleSyncBlock: Failed to get last block while syncing block hash=%s height=%d: %v", block.Hash, block.Height, err)
			return err
		}

		exists, err := j.dbChain.ExistingBlock(ctx, block.Hash, tx)
		if err != nil {
			log.Errorf("handleSyncBlock: Failed to check existence of block hash=%s height=%d: %v", block.Hash, block.Height, err)
			return err
		}

		if exists {
			log.Infof("handleSyncBlock: Block hash=%s height=%d already exists, skipping", block.Hash, block.Height)
			continue
		}

		if j.isGenesisBlock(block) {
			log.Infof("handleSyncBlock: Processing genesis block hash=%s", block.Hash)
			err := j.handleCreateBlock(block, tx)
			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create genesis block hash=%s: %v", block.Hash, err)
				return err
			}
			err = j.handleCreateTransactions(ctx, block.Transactions, block, tx)

			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create transactions for genesis block hash=%s: %v", block.Hash, err)
				return err
			}

			err = j.handleCreateUtxo(block, tx)

			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create UTXOs for genesis block hash=%s: %v", block.Hash, err)
				return err
			}
			log.Infof("handleSyncBlock: Completed processing genesis block hash=%s", block.Hash)
			continue
		} else {
			existsPrevBlock, err := j.dbChain.ExistingBlock(ctx, block.PrevHash, tx)
			if err != nil {
				log.Errorf("handleSyncBlock: Failed to check prev block hash=%s existence for block hash=%s height=%d: %v", block.PrevHash, block.Hash, block.Height, err)
				return err
			}

			if !existsPrevBlock {
				log.Warnf("handleSyncBlock: PrevHash=%s not exist for block hash=%s height=%d, skipping sync", block.PrevHash, block.Hash, block.Height)
				return nil
			}
			log.Infof("handleSyncBlock: Prev block hash=%s exists for block hash=%s height=%d", block.PrevHash, block.Hash, block.Height)
		}

		lastBlockNChainWork := big.NewInt(0)
		lastBlockNChainWork.SetString(lastBlock.NchainWork, 10)

		if block.NChainWork.Cmp(lastBlockNChainWork) > 0 {
			log.Infof("handleSyncBlock: Block hash=%s height=%d has higher chain work than last block (%s > %s), proceeding to add", block.Hash, block.Height, block.NChainWork, lastBlock.NchainWork)
			_, err = j.dbChain.GetBlockByHash(ctx, block.PrevHash, tx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Infof("handleSyncBlock: Prev block hash=%s not found in DB for block hash=%s height=%d, initiating chain switch", block.PrevHash, block.Hash, block.Height)
					if err := j.handleSwitchChain(block, tx); err != nil {
						log.Errorf("handleSyncBlock: Failed chain switch for block hash=%s height=%d: %v", block.Hash, block.Height, err)
						return err
					}
					continue
				}
				log.Errorf("handleSyncBlock: Failed to get prev block hash=%s for block hash=%s height=%d: %v", block.PrevHash, block.Hash, block.Height, err)
				return err
			}
			log.Debugf("handleSyncBlock: Confirmed prev block hash=%s exists for block hash=%s height=%d", block.PrevHash, block.Hash, block.Height)

			err := j.handleCreateBlock(block, tx)
			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create block hash=%s height=%d: %v", block.Hash, block.Height, err)
				return err
			}

			err = j.handleCreateTransactions(ctx, block.Transactions, block, tx)

			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create transactions for block hash=%s height=%d: %v", block.Hash, block.Height, err)
				return err
			}

			err = j.handleCreateUtxo(block, tx)

			if err != nil {
				log.Errorf("handleSyncBlock: Failed to create UTXOs for block hash=%s height=%d: %v", block.Hash, block.Height, err)
				return err
			}
			log.Infof("handleSyncBlock: Successfully added block hash=%s height=%d to chain", block.Hash, block.Height)

		} else {
			log.Warnf("handleSyncBlock: Block hash=%s height=%d chain work not higher than last (%s <= %s), skipping", block.Hash, block.Height, block.NChainWork, lastBlock.NchainWork)
		}
	}

	log.Infof("handleSyncBlock: Completed sync for %d blocks", len(blocks))
	return nil
}

func (j *jobBlockSync) StartBlockSync(interval time.Duration) {
	log.Info("🔄 Chain sync started")
	ctx := context.Background()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {

		lastBlock, err := j.dbChain.GetLastBlock(ctx, nil)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			lastBlock = dbchain.Block{
				Height: 1,
			}
			log.Warnf("StartBlockSync: No last block in DB, starting from height=1")
		} else if err != nil {
			log.Errorf("❌ StartBlockSync: GetLastBlock error: %v", err)
			continue
		}

		log.Infof("⏱️ StartBlockSync: Sync tick - current chain height=%d, fetching from height=%d", lastBlock.Height, lastBlock.Height+1)

		blocks, err := j.dbChainRpc.GetBlocksByHeightRange(lastBlock.Height, 10)
		if err != nil {
			log.Errorf("❌ StartBlockSync: RPC call error for height range starting %d: %v", lastBlock.Height, err)
			continue
		}

		if len(blocks) == 0 {
			log.Warn("⚠️ StartBlockSync: No new blocks received from RPC")
			continue
		}

		log.Infof("📦 StartBlockSync: Received %d new blocks from RPC (heights %d to %d)", len(blocks), blocks[0].Height, blocks[len(blocks)-1].Height)

		tx, err := db.Psql.BeginTx(ctx, nil)
		if err != nil {
			log.Errorf("❌ StartBlockSync: BeginTx error: %v", err)
			continue
		}
		log.Debugf("StartBlockSync: Started transaction for syncing %d blocks", len(blocks))

		err = j.handleSyncBlock(blocks, tx)
		if err != nil {
			log.Errorf("❌ StartBlockSync: Sync block error for %d blocks: %v", len(blocks), err)
			if err := tx.Rollback(); err != nil {
				log.Errorf("❌ StartBlockSync: Rollback error after sync failure: %v", err)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Errorf("❌ StartBlockSync: Commit transaction error after sync: %v", err)
			continue
		}

		log.Infof("✅ StartBlockSync: Sync complete - successfully added up to height=%d (%d blocks processed)", blocks[0].Height, len(blocks))

	}

}
