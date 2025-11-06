package chain

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/utils"
	dbchain "ChainServer/internal/db/chain"
	"ChainServer/internal/states"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (s *ChainService) GetBlocks(dto dto.PaginationQuery) ([]dbchain.GetListBlocksRow, *response.PaginationMeta, *apperror.AppError) {
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
			return make([]dbchain.GetListBlocksRow, 0), nil, nil
		}
		return nil, nil, apperror.Internal("Failted to get blocks", err)
	}

	lastestBlock, err := s.dbRepo.GetLastBlock(ctx, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]dbchain.GetListBlocksRow, 0), nil, nil
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

func (s *ChainService) GetBlockDetail(dto *GetBlockDetailDto) (BlockDetail, *apperror.AppError) {
	ctx := context.Background()

	limit := int32(*dto.Limit)
	page := int32(*dto.Page)

	offsetTx := limit * (page - 1)

	block, err := s.dbRepo.GetBlockDetailWithTransactions(ctx, dbchain.GetBlockDetailWithTransactionsParams{
		BID:      dto.BlockHash,
		OffsetTx: offsetTx,
		LimitTx:  limit,
	})

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return BlockDetail{}, apperror.Internal("Failed to get block detail", err)
		}
		return BlockDetail{}, apperror.NotFound("Block not found", nil)
	}

	pagination := helpers.BuildPaginationMeta(
		int64(limit),
		int64(page),
		block.TxCount,
		nil,
	)

	difficulty, _ := utils.CompactToDifficulty(uint32(block.Nbits)).Int64()

	blockDetail := BlockDetail{
		ID:         block.ID,
		BID:        block.BID,
		PrevHash:   block.PrevHash,
		Nonce:      block.Nonce,
		Height:     block.Height,
		MerkleRoot: block.MerkleRoot,
		Nbits:      block.Nbits,
		TxCount:    block.TxCount,
		NchainWork: block.NchainWork,
		Size:       block.Size,
		Timestamp:  block.Timestamp,
		Difficulty: difficulty,
		Miner:      block.Miner,
		TotalFee:   block.TotalFee,
		Transactions: struct {
			Data json.RawMessage
			Meta *response.PaginationMeta
		}{
			Data: block.Transactions,
			Meta: pagination,
		},
	}

	return blockDetail, nil
}

func (s *ChainService) GetNetwork() (*NetworkInfo, *apperror.AppError) {
	ctx := context.Background()

	bestHeight, err := s.dbRepo.GetBestHeight(ctx, nil)
	if err != nil {
		log.Errorf("Failed to get best height: %v", err)
		return nil, apperror.Internal("Something went wrong, please try again", nil)
	}

	recentBlocks, err := s.dbRepo.GetRecentBlocksForNetworkInfo(ctx, 10)
	if err != nil {
		log.Errorf("Failed to get recent blocks: %v", err)
		return nil, apperror.Internal("Something went wrong, please try again", nil)
	}

	var nbitsList []uint32
	var timestamps []int64
	for _, block := range recentBlocks {
		nbitsList = append(nbitsList, uint32(block.Nbits))
		timestamps = append(timestamps, block.Timestamp)
	}

	avgDifficulty := utils.AverageDifficulty(nbitsList)
	avgBlockTimes := utils.AverageBlockTime(timestamps)
	hashrate := utils.CalculateHashrate(avgDifficulty, avgBlockTimes)
	_, unit := utils.FormatHashrate(hashrate)

	count, err := s.tranRepo.CountPendingTxs(ctx, nil)
	if err != nil {
		log.Errorf("Failed to get count tx pending: %v", err)
		return nil, apperror.Internal("Something went wrong, please try again", nil)
	}

	return &NetworkInfo{
		LastBlock:     bestHeight,
		Hashrate:      fmt.Sprintf("%.2f %s", hashrate, unit),
		AvgBlockTime:  avgBlockTimes,
		AvgDifficulty: avgDifficulty,
		SyncStatus:    states.ChainSyncState.SyncStatus,
		TxPending:     count,
		NetworkHealth: utils.EvaluateNetwork(avgBlockTimes, float64(10*time.Minute)),
	}, nil
}

func (s *ChainService) GetMiners(dto *dto.PaginationQuery) ([]dbchain.GetMinersRow, *response.PaginationMeta, *apperror.AppError) {
	ctx := context.Background()

	limit := int32(*dto.Limit)
	offset := int32((*dto.Page - 1)) * limit

	arg := dbchain.GetMinersParams{
		Offset: offset,
		Limit:  limit,
	}

	miners, err := s.dbRepo.GetMiners(ctx, arg)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperror.Internal("Something went wrong, please try again!", nil)
		}
		miners = make([]dbchain.GetMinersRow, 0)
	}

	count, err := s.dbRepo.CountMiners(ctx)
	if err != nil {
		return nil, nil, apperror.Internal("Something went wrong, please try again!", nil)
	}

	pagination := helpers.BuildPaginationMeta(
		int64(limit),
		*dto.Page,
		count,
		nil,
	)

	return miners, pagination, nil
}
