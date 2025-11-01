package transaction

import (
	"ChainServer/internal/app/module/utxo"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/types"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/db"
	dbchain "ChainServer/internal/db/chain"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	dbutxo "ChainServer/internal/db/utxo"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TransactionService struct {
	dbRepo   DbTransactionRepository
	utxoRepo utxo.DbUTXORepository
}

func NewTransactionService(
	dbRepo DbTransactionRepository,
	utxoRepo utxo.DbUTXORepository,
) *TransactionService {
	return &TransactionService{
		dbRepo:   dbRepo,
		utxoRepo: utxoRepo,
	}
}

func (s *TransactionService) GetListTransaction(dto dto.PaginationQuery) ([]dbchain.Transaction, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	page := *dto.Page
	limit := *dto.Limit

	txs, err := s.dbRepo.GetListTransaction(ctx, dbchain.GetListTransactionsParams{
		Offset: int32((page - 1)) * int32(limit),
		Limit:  int32(limit),
	}, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.Transaction, 0), nil, nil
		}

		return nil, nil, apperror.Internal("Failed to get transactions", err)
	}

	count, err := s.dbRepo.GetCountTransaction(ctx)

	if err != nil {
		return nil, nil, apperror.Internal("Failed to get count transaction", err)
	}

	pagination := helpers.BuildPaginationMeta(
		limit,
		page,
		count,
		dto.NextCursor,
	)

	return txs, pagination, nil
}

func (s *TransactionService) CreateNewTransaction(payload *utils.JWTPayload[types.JWTWalletAuthPayload], dto *NewTransactionParsed) (*string, *apperror.AppError) {
	ctx := context.Background()

	internalErrCommon := apperror.Internal("Something went wrong. Please try again.", nil)

	pubKeyBytes, err := hex.DecodeString(payload.Data.Pubkey)
	if err != nil {
		return nil, internalErrCommon
	}

	pubKeyHash := utils.PublicKeyHash(pubKeyBytes)

	utxos, err := s.utxoRepo.FindUTXOs(ctx, hex.EncodeToString(pubKeyHash), nil)

	if err != nil {
		return nil, internalErrCommon
	}

	acc := utils.NewCoinAmountFromFloat(0.0)
	var spendable []dbutxo.Utxo

	prevTxs := map[string]dbutxo.Utxo{}

	for _, utxo := range utxos {
		value, err := strconv.ParseFloat(utxo.Value, 64)
		if err != nil {
			log.Errorf("Failed To Parse utxo: %v", err)
			return nil, internalErrCommon
		}
		valueTx := utils.NewCoinAmountFromFloat(value)
		acc = acc.Add(valueTx)
		spendable = append(spendable, utxo)
		prevTxs[utxo.TxID] = utxo

		amount := utils.NewCoinAmountFromFloat(dto.Data.Amount)
		fee := utils.NewCoinAmountFromFloat(dto.Data.Fee)

		if acc.Cmp(amount.Add(fee)) >= 0 {
			break
		}
	}

	fromAddrByte, err := utils.Base58Decode(strings.Trim(payload.Data.Address, " "))
	if err != nil {
		return nil, internalErrCommon
	}

	tx, apperr := NewTransaction(
		pubKeyBytes,
		fromAddrByte,
		dto.Data.To,
		dto.Data.Amount,
		dto.Data.Fee,
		spendable,
		acc.ToFloat(),
	)

	if apperr != nil {
		return nil, apperr
	}

	txWithSigning, apperr := tx.WithSigning(prevTxs)
	if apperr != nil {
		return nil, apperr
	}

	data, err := json.Marshal(txWithSigning)
	if err != nil {
		return nil, internalErrCommon
	}

	txEncode, err := utils.EncryptData(data, env.Cfg.Encode_data_secret_Key)
	if err != nil {
		log.Errorf("EncryptData error: %v", err)
		return nil, internalErrCommon
	}

	return &txEncode, nil
}

func (s *TransactionService) SendTransaction(payload *utils.JWTPayload[types.JWTWalletAuthPayload], dto *SendTransactionDataParsed) *apperror.AppError {
	ctx := context.Background()

	internalServerErr := apperror.Internal("Transaction processing failled. Please try again later.", nil)

	pubKeyBytes, err := hex.DecodeString(payload.Data.Pubkey)
	if err != nil {
		return apperror.BadRequest("Invalid wallet public key format.", err)
	}

	pubkeyHashByte := utils.PublicKeyHash(pubKeyBytes)
	pubKeyHash := hex.EncodeToString(pubkeyHashByte)

	utoxs, err := s.utxoRepo.FindUTXOs(ctx, pubKeyHash, nil)
	if err != nil {
		log.Errorf("FindUTXO error: %v", err)
		return internalServerErr
	}

	if len(utoxs) == 0 {
		return apperror.BadRequest("No enough funds in your wallet", nil)
	}

	existingPendingTx, err := s.dbRepo.ExistingPendingTransaction(
		ctx,
		dbPendingTx.PendingTxExistsParams{
			TxID:   hex.EncodeToString(dto.Transaction.ID),
			Status: []string{string(constants.TxStatusPending), string(constants.TxStatusMining)},
		},
		nil)
	if err != nil {
		log.Errorf("Check existing pending transaction with error: %v", err)
		return internalServerErr
	}

	if existingPendingTx {
		return apperror.BadRequest("Transaction already exists", nil)
	}

	prevTxs := map[string]dbutxo.Utxo{}

	for _, tx := range utoxs {
		prevTxs[tx.TxID] = tx
	}

	sigOk := VerifyTransactionSig(&dto.Transaction, prevTxs)

	if !sigOk {
		return apperror.BadRequest("Transaction signature verification failed.", nil)
	}

	_, ok := constants.Priorities[dto.Priority]
	if !ok {
		return apperror.BadRequest("Priority level is not valid", nil)
	}

	tx, err := db.Psql.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Errorf("BeginTx error: %v", err)
		return internalServerErr
	}

	payloadNewPendingTx := dbPendingTx.InsertPendingTransactionParams{
		TxID:            hex.EncodeToString(dto.Transaction.ID),
		Address:         payload.Data.Address,
		Status:          string(constants.TxStatusPending),
		Priority:        helpers.Int32ToNullInt32(int32(dto.Priority)),
		ReceiverAddress: dto.ReceiverAddr,
		Message:         helpers.StringToNullString(dto.Message),
		Amount:          fmt.Sprintf("%.8f", dto.Amount),
		Fee:             fmt.Sprintf("%.8f", dto.Fee),
	}

	newTxPending, err := s.dbRepo.InsertPendingTransaction(ctx, payloadNewPendingTx, tx)
	if err != nil {
		log.Errorf("Insert new pending transaction error: %v", err)

		tx.Rollback()
		return internalServerErr
	}

	txEncode, err := json.Marshal(dto.Transaction)
	if err != nil {
		log.Errorf("Marshal transaction error: %v", err)
		tx.Rollback()
		return internalServerErr
	}

	payloadNewPendingTxData := dbPendingTx.InsertPendingTxDataParams{
		TxRef:      newTxPending.ID,
		PubKeyHash: pubKeyHash,
		RawTx:      txEncode,
	}

	_, err = s.dbRepo.InsertPendingTxData(ctx, payloadNewPendingTxData, tx)
	if err != nil {
		log.Errorf("Insert new pending transaction data error: %v", err)
		tx.Rollback()
		return internalServerErr
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("Commit transaction error: %v", err)
		tx.Rollback()
		return internalServerErr
	}

	return nil
}

func (s *TransactionService) TransactionPending(payload *utils.JWTPayload[types.JWTWalletAuthPayload], pagination *dto.PaginationQuery) ([]dbPendingTx.PendingTxsByAddressAndStatusRow, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	internalErrCommon := apperror.Internal("Something went wrong. Please try again.", nil)

	limit := *pagination.Limit
	page := *pagination.Page

	txPending, err := s.dbRepo.FindTxPendingByAddrAndStatus(ctx, dbPendingTx.PendingTxsByAddressAndStatusParams{
		Address: payload.Data.Address,
		Limit:   int32(limit),
		Offset:  (int32(page) - 1) * int32(limit),
		Status: []string{
			string(constants.TxStatusPending),
			string(constants.TxStatusMining),
		},
	}, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, internalErrCommon
	}

	count, err := s.dbRepo.CountPendingTxByAddr(ctx, payload.Data.Address, nil)
	if err != nil {
		return nil, nil, internalErrCommon
	}

	paginationMeta := helpers.BuildPaginationMeta(
		limit,
		page,
		count,
		pagination.NextCursor,
	)

	return txPending, paginationMeta, nil
}

func (s *TransactionService) SearchTransactions(queries *GetTransactionSearchDto) ([]dbchain.Transaction, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	limit := int32(*queries.Limit)
	page := int32(*queries.Page)

	if queries.Search_Tx_Query == "" || len(queries.Search_Tx_Query) < 1 {
		txs, err := s.dbRepo.GetListTransactionByBlockHash(ctx, dbchain.GetListTransactionByBIDParams{
			BID:    queries.B_Hash,
			Offset: (page - 1) * limit,
			Limit:  limit,
		}, nil)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, apperror.NotFound("No transactions found matching the search criteria", nil)
			}
			log.Errorf("Get list transactions error: %v", err)
			return nil, nil, apperror.Internal("Failed to search transactions", nil)
		}

		count, err := s.dbRepo.CountTransactionByBID(ctx, queries.B_Hash, nil)
		if err != nil {
			log.Errorf("Count transactions by block hash error: %v", err)
			return nil, nil, apperror.Internal("Failed to get count of searched transactions", nil)
		}

		pagination := helpers.BuildPaginationMeta(
			int64(limit),
			int64(page),
			count,
			queries.NextCursor,
		)

		return txs, pagination, nil
	}

	txs, err := s.dbRepo.SearchFuzzyTransactionsByBlock(ctx, dbchain.SearchFuzzyTransactionsByBlockParams{
		SearchQuery: queries.Search_Tx_Query,
		BHash:       queries.B_Hash,
		Limit:       limit,
		Offset:      (page - 1) * limit,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperror.NotFound("No transactions found matching the search criteria", nil)
		}
		log.Errorf("SearchFuzzyTransactions error: %v", err)
		return nil, nil, apperror.Internal("Failed to search transactions", nil)
	}

	if len(txs) == 0 {
		return nil, nil, apperror.NotFound("No transactions found matching the search criteria", nil)
	}

	count, err := s.dbRepo.CountFuzzyTransactionsByBlock(ctx, dbchain.CountFuzzyTransactionsByBlockParams{
		SearchQuery: queries.Search_Tx_Query,
		BHash:       queries.B_Hash,
	})

	if err != nil {
		log.Errorf("CountFuzzyTransactions error: %v", err)
		return nil, nil, apperror.Internal("Failed to get count of searched transactions", nil)
	}

	pagination := helpers.BuildPaginationMeta(
		int64(limit),
		int64(page),
		count,
		queries.NextCursor,
	)

	return txs, pagination, nil
}

func (s *TransactionService) GetPendingTransactions(queries *GetTransactionPendingDto) ([]dbPendingTx.GetListPendingTxsRow, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	limit := int32(*queries.Limit)
	page := int32(*queries.Page)

	txs, err := s.dbRepo.GetPendingTransactions(ctx, dbPendingTx.GetListPendingTxsParams{
		Offset: (page - 1) * limit,
		Limit:  limit,
	})

	if err != nil {
		log.Errorf("Failed to get pending transactions: %v", err)
		return nil, nil, apperror.Internal("Internal server", nil)
	}

	count, err := s.dbRepo.GetCountPendingTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to get count pending transactions: %v", err)
		return nil, nil, apperror.Internal("Internal server", nil)
	}

	pagination := helpers.BuildPaginationMeta(
		int64(limit),
		int64(page),
		count,
		nil,
	)

	return txs, pagination, nil
}
