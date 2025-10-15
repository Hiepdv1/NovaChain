package dashboard

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/utils"
	dbchain "ChainServer/internal/db/chain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

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

	// --- fetch first and last block ---
	firstBlock, err := s.chainRepo.GetBlockByHeight(ctx, 1, nil)
	if err != nil {
		log.Errorf("[Dashboard] ❌ Failed to fetch first block (height=1): %v", err)
		return nil, apperror.Internal("Something went wrong. Please try again.", err)
	}

	lastBlock, err := s.chainRepo.GetBlockByHeight(ctx, bestHeight, nil)
	if err != nil {
		log.Errorf("[Dashboard] ❌ Failed to fetch last block (height=%d): %v", bestHeight, err)
		return nil, apperror.Internal("Something went wrong. Please try again.", err)
	}

	// compute average block time per block (safely)
	var avgBlockTimePerBlock float64
	if lastBlock.Height > firstBlock.Height {
		totalTime := lastBlock.Timestamp - firstBlock.Timestamp
		if totalTime <= 0 {
			log.Warn("[Dashboard] ⚠️ Non-positive total time between first and last block — cannot compute avg block time")
			avgBlockTimePerBlock = 0
		} else {
			avgBlockTimePerBlock = float64(totalTime) / float64(lastBlock.Height-firstBlock.Height)
		}
	} else {
		// not enough height span to compute meaningful average
		log.Warn("[Dashboard] ⚠️ Not enough height span to compute avgBlockTime (need last.height > first.height)")
		avgBlockTimePerBlock = 0
	}

	// difficulty -> hashrate (H/s). Use 2^32 constant instead of math.Pow each time.
	const two32 = 4294967296.0
	difficultyVal, _ := utils.CompactToDifficulty(uint32(lastBlock.Nbits)).Float64()
	var hCurrent float64
	if avgBlockTimePerBlock > 0 {
		hCurrent = (difficultyVal * two32) / avgBlockTimePerBlock
	} else {
		hCurrent = 0
	}

	// --- collect 24h blocks and compute averages ---
	listBlock, err := s.chainRepo.GetListBlockByHours(ctx, 24, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("[Dashboard] ⚠️ No block data found for the last 24 hours")
			// continue: we will return overview with zeros / N/A
		} else {
			log.Errorf("[Dashboard] ❌ GetListBlockByHours failed: %v", err)
			return nil, apperror.Internal("Something went wrong. Please try again.", err)
		}
	}

	var avgBlockTime24h float64
	var avgDifficulty24h float64
	var hashrate24h float64
	var changeStr string = "N/A"

	if len(listBlock) >= 2 {
		N := len(listBlock)
		intervals := N - 1
		firstTs := listBlock[intervals].Timestamp
		lastTs := listBlock[0].Timestamp

		totalTime := lastTs - firstTs
		if totalTime > 0 {
			avgBlockTime24h = float64(totalTime) / float64(intervals)

			sumDiff := 0.0
			for _, b := range listBlock {
				d, _ := utils.CompactToDifficulty(uint32(b.Nbits)).Float64()
				sumDiff += d
			}
			avgDifficulty24h = sumDiff / float64(N)

			if avgBlockTime24h > 0 {
				hashrate24h = (avgDifficulty24h * two32) / avgBlockTime24h
			} else {
				hashrate24h = 0
			}

			const eps = 1e-12
			if hashrate24h <= eps {
				if hCurrent <= eps {
					changeStr = "0.00%"
				} else {
					changeStr = "0.00%"
				}
			} else {
				change := ((hCurrent - hashrate24h) / hashrate24h) * 100
				if math.IsNaN(change) || math.IsInf(change, 0) {
					changeStr = "0.00%"
				} else {
					changeStr = fmt.Sprintf("%.2f%%", change)
				}
			}
		} else {
			log.Warn("[Dashboard] ⚠️ Non-positive total time within 24h blocks — cannot compute 24h averages")
			avgBlockTime24h = 0
			hashrate24h = 0
			changeStr = "0.00%"
		}
	} else {
		log.Warnf("[Dashboard] ⚠️ Not enough blocks (%d) to calculate 24h averages", len(listBlock))
		changeStr = "0.00%"
	}

	totalTransaction, err := s.tranRepo.CountTransactions(ctx, nil)
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

	value, unit := utils.FormatHashrate(hCurrent)
	log.Infof("[Dashboard] ✅ Network overview collected successfully (height=%d, hashrate=%.3f %s)", bestHeight, value, unit)

	return &NetworkOverview{
		Chain: struct {
			BestHeight int64
			PerHours   int64
		}{
			BestHeight: bestHeight,
			PerHours:   blockCountByHours,
		},
		Hashrate: struct {
			Value  string
			Per24H string
		}{
			Value:  fmt.Sprintf("%.2f %s", value, unit),
			Per24H: changeStr,
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
		listBlocks = []dbchain.Block{}
	}

	if listTxs == nil {
		listTxs = []dbchain.Transaction{}
	}

	return &RecentActivity{
		Blocks: listBlocks,
		Txs:    listTxs,
	}, nil
}
