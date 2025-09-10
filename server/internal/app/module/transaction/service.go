package transaction

import (
	"ChainServer/internal/app/module/utxo"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/types"
	"ChainServer/internal/common/utils"
	dbchain "ChainServer/internal/db/chain"
	dbutxo "ChainServer/internal/db/utxo"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
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

	page := int32(*dto.Page)
	limit := int32(*dto.Limit)

	txs, err := s.dbRepo.GetListTransaction(ctx, dbchain.GetListTransactionsParams{
		Offset: (page - 1) * limit,
		Limit:  *dto.Limit,
	})

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
		int32(count),
		dto.NextCursor,
	)

	return txs, pagination, nil
}

func (s *TransactionService) CreateNewTransaction(payload *utils.JWTPayload[types.JWTWalletAuthPayload], dto *NewTransactionParsed) (any, *apperror.AppError) {
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

	var acc float64
	var spendable []dbutxo.Utxo

	for _, utxo := range utxos {
		value, err := strconv.ParseFloat(utxo.Value, 64)
		if err != nil {
			log.Errorf("Failed To Parse utxo: %v", err)
			return nil, internalErrCommon
		}
		acc += value
		spendable = append(spendable, utxo)
		if acc >= dto.Data.Amount+dto.Data.Fee {
			break
		}
	}

	fromAddrByte, err := utils.Base58Decode(strings.Trim(payload.Data.Address, " "))
	if err != nil {
		return nil, internalErrCommon
	}

	tx, apperr := utils.NewTransaction(
		dto.PubKey,
		fromAddrByte,
		dto.Data.To,
		dto.Data.Amount,
		dto.Data.Fee,
		spendable,
		acc,
	)

	if apperr != nil {
		return nil, apperr
	}

	return tx, nil
}
