package dashboard

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/utils"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type DashboardService struct {
	tranRepo  TXRepository
	chainRepo ChainRepository
}

func NewDashboardService(
	tranRepo TXRepository,
	chainRepo ChainRepository,
) *DashboardService {
	return &DashboardService{
		tranRepo:  tranRepo,
		chainRepo: chainRepo,
	}
}

// Updated GetNetworkOverview — paste this function into your DashboardService file (replace the old one).
// Improvements:
// - compute average block time per-block (divide by height difference)
// - avoid div-by-zero and negative timestamps
// - compute 24h avg difficulty dividing by number of blocks (not intervals)
// - safe percent-change when historical hashrate is zero/near-zero (show "N/A" or "0.00%")
// - adaptive hashrate unit formatting (H/s, kH/s, MH/s, GH/s, TH/s)

func (s *DashboardService) GetNetworkOverview() (*NetworkOverview, *apperror.AppError) {
	ctx := context.Background()

	bestHeight, err := s.chainRepo.GetBestHeight(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ GetBestHeight failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No blocks found in GetBestHeight (setting height = 0)")
		bestHeight = 0
	}

	blockCountByHours, err := s.chainRepo.GetBlockCountByHours(ctx, 1, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ GetBlockCountByHours failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No block data for the last hour (PerHours = 0)")
		blockCountByHours = 0
	}

	recentBlocks, err := s.chainRepo.GetRecentBlocksForNetworkInfo(ctx, 200)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ Get recent blocks failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No block found in recentBlocks")
	}

	currentBlocks := recentBlocks[:100]
	previousBlocks := recentBlocks[100:]

	var nbitsList []uint32
	var timestamps []int64
	for _, block := range currentBlocks {
		nbitsList = append(nbitsList, uint32(block.Nbits))
		timestamps = append(timestamps, block.Timestamp)
	}

	avgDifficulty := utils.AverageDifficulty(nbitsList)
	avgBlockTimes := utils.AverageBlockTime(timestamps)
	hashrate := utils.CalculateHashrate(avgDifficulty, avgBlockTimes)
	hashrateFormat, unit := utils.FormatHashrate(hashrate)

	var prevNbitsList []uint32
	var prevTimestamp []int64
	for _, block := range previousBlocks {
		prevNbitsList = append(prevNbitsList, uint32(block.Nbits))
		prevTimestamp = append(prevTimestamp, block.Timestamp)
	}
	prevAvgDifficulty := utils.AverageDifficulty(prevNbitsList)
	prevAvgBlockTimes := utils.AverageBlockTime(prevTimestamp)
	prevHashrate := utils.CalculateHashrate(prevAvgDifficulty, prevAvgBlockTimes)

	hashrateChange := utils.CalculateHashrateChange(hashrate, prevHashrate)

	totalTransaction, err := s.tranRepo.GetCountTransaction(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ CountTransactions failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No transactions found in CountTransactions")
	}

	countTodayTx, err := s.tranRepo.CountTodayTransaction(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ CountTodayTransaction failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No transactions found for today")
	}

	totalPendingTxs, err := s.tranRepo.CountPendingTxs(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ CountPendingTxs failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No pending transactions found")
	}

	countTodayPendingTxs, err := s.tranRepo.CountTodayPendingTxs(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ CountTodayPendingTxs failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No pending transactions added today")
	}

	countMinerWork, err := s.chainRepo.CountDistinctMiners(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ CountDistinctMiners failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No miner data found")
	}

	countTodayMinerWork, err := s.chainRepo.GetCountTodayWorkerMiners(ctx, nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[Dashboard] ❌ GetCountTodayWorkerMiners failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
		log.Warn("[Dashboard] ⚠️ No miner activity found today")
	}

	return &NetworkOverview{
		Chain: struct {
			BestHeight int64
			PerHours   int64
		}{
			BestHeight: bestHeight,
			PerHours:   blockCountByHours,
		},
		Hashrate: struct {
			Value      string
			Trend      string
			ChangeRate string
		}{
			Value:      fmt.Sprintf("%.2f %s", hashrateFormat, unit),
			ChangeRate: fmt.Sprintf("%.2f%%", hashrateChange),
			Trend:      utils.FormatTrend(hashrateChange),
		},
		Transaction: struct {
			Total      int64
			AddedToday int64
		}{
			Total:      totalTransaction,
			AddedToday: countTodayTx,
		},
		PendingTx: struct {
			Count      int64
			AddedToday int64
		}{
			Count:      totalPendingTxs,
			AddedToday: countTodayPendingTxs,
		},
		ActiveMiners: struct {
			Count  int64
			Worker int64
		}{
			Count:  countMinerWork,
			Worker: countTodayMinerWork,
		},
	}, nil
}

func (s *DashboardService) GetRecentActivity() (*RecentActivity, *apperror.AppError) {
	ctx := context.Background()

	listBlocks, err := s.chainRepo.GetListBlock(ctx, dbchain.GetListBlocksParams{
		Offset: 0,
		Limit:  5,
	}, nil)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.Internal("Something went wrong, please try again.", nil)
		}
	}

	listTxs, err := s.tranRepo.GetListTransaction(ctx, dbchain.GetListTransactionsParams{
		Offset: 0,
		Limit:  5,
	}, nil)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.Internal("Something went wrong, please try again.", nil)
		}
	}

	if listBlocks == nil {
		listBlocks = []dbchain.GetListBlocksRow{}
	}

	if listTxs == nil {
		listTxs = []dbchain.Transaction{}
	}

	return &RecentActivity{
		Blocks: listBlocks,
		Txs:    listTxs,
	}, nil
}
