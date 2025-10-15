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
	"strings"

	log "github.com/sirupsen/logrus"
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

	limit := *dto.Limit
	page := *dto.Page
	offset := int32((page - 1)) * int32(limit)

	blocks, err := s.dbRepo.GetListBlock(ctx, dbchain.GetListBlocksParams{
		Offset: offset,
		Limit:  int32(limit),
	}, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.Block, 0), nil, nil
		}
		return nil, nil, apperror.Internal("Failted to get blocks", err)
	}

	lastestBlock, err := s.dbRepo.GetLastBlock(ctx, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.Block, 0), nil, nil
		}
		return nil, nil, apperror.Internal("Failted to get blocks", err)
	}

	pagination := helpers.BuildPaginationMeta(
		limit,
		page,
		lastestBlock.Height,
		nil,
	)

	return blocks, pagination, nil
}

func (s *ChainService) GetSearchResult(dto *GetSearchResultDto) ([]dbchain.SearchFuzzyRow, *response.PaginationMeta, *apperror.AppError) {

	ctx := context.Background()

	limit := int32(*dto.Limit)
	page := int32(*dto.Page)

	offset := limit * (page - 1)

	searchQuery := strings.TrimSpace(dto.Search_Query)

	total, err := s.dbRepo.CountFuzzy(ctx, searchQuery)
	if err != nil {
		log.Errorf("[Search-Fizzy] Content %s error %v", dto.Search_Query, err)
		return nil, nil, apperror.Internal("Something went wrong. Please try again!", nil)
	}

	resultFizzy, err := s.dbRepo.SearchFuzzy(ctx, dbchain.SearchFuzzyParams{
		SearchQuery: searchQuery,
		Offset:      offset,
		Limit:       limit,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperror.NotFound("Not found", nil)
		}

		log.Errorf("[Search-Fizzy] Content %s error %v", dto.Search_Query, err)
		return nil, nil, apperror.Internal("Something went wrong. Please try again!", nil)

	}

	pagination := helpers.BuildPaginationMeta(
		int64(limit),
		int64(page),
		total,
		nil,
	)

	return resultFizzy, pagination, nil

}
