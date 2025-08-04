package chain

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
	"errors"
)

type ChainService struct {
	rpcRepo  RPCChainRepository
	dbRepo   DBChainRepository
	tranRepo transaction.DbTransactionRepository
}

func NewChainService(
	rpcRepo RPCChainRepository,
	dbRepo DBChainRepository,
	tranRepo transaction.DbTransactionRepository,
) *ChainService {
	return &ChainService{
		rpcRepo:  rpcRepo,
		dbRepo:   dbRepo,
		tranRepo: tranRepo,
	}
}

func (s *ChainService) GetBlocks(dto dto.PaginationQuery) ([]dbchain.Block, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	limit := int32(*dto.Limit)
	page := int32(*dto.Page)
	offset := (page - 1) * limit

	blocks, err := s.dbRepo.GetListBlock(ctx, dbchain.GetListBlocksParams{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.Block, 0), nil, nil
		}
		return nil, nil, apperror.Internal("Failted to get blocks", err)
	}

	lastestBlock, err := s.dbRepo.GetLastBlock(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.Block, 0), nil, nil
		}
		return nil, nil, apperror.Internal("Failted to get blocks", err)
	}

	pagination := helpers.BuildPaginationMeta(
		limit,
		page,
		int32(lastestBlock.Height),
		nil,
	)

	return blocks, pagination, nil
}
