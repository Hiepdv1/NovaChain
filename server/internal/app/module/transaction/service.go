package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
	"errors"
)

type TransactionService struct {
	dbRepo DbTransactionRepository
}

func NewTransactionService(
	dbRepo DbTransactionRepository,
) *TransactionService {
	return &TransactionService{
		dbRepo: dbRepo,
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
