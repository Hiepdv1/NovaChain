package p2p

import (
	"context"
	blockchain "core-blockchain/core"
	"core-blockchain/memopool"
	"maps"
	"slices"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (net *Network) MinersEventLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go net.requestTransaction(ctx)
	net.miningLoop(ctx)
}

func (net *Network) requestTransaction(ctx context.Context) {
	poolCheckTicker := time.NewTicker(30 * time.Second)
	defer poolCheckTicker.Stop()
	for range poolCheckTicker.C {
		select {
		case <-ctx.Done():
			return
		default:
			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{
					net.Host.ID().String(),
				},
				func(p peer.ID) {
					net.SendRequestTxFromPool(
						p.String(),
						TxFromPool{
							SendFrom: net.Host.ID().String(),
							Count:    1,
						},
					)
				},
			)
		}
	}
}

func (net *Network) miningLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			net.tryMine(ctx)
		}
	}
}

func (net *Network) tryMine(ctx context.Context) {

	if len(MemoryPool.Pending) == 0 {
		return
	}

	txs := MemoryPool.SelectHighFeeTx()

	if len(MemoryPool.Queued) == 0 {
		return
	}

	if !net.syncCompleted {
		return
	}

	log.Info("Starting mining new block...")
	net.mineBlock(ctx, txs)
}

func (net *Network) mineBlock(ctx context.Context, txs map[string]blockchain.Transaction) {
	ctxInternal, cancel := context.WithCancel(context.Background())
	defer cancel()

	chain := net.Blockchain
	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		log.Errorf("Get last block with err: %v", err)
		return
	}

	log.Infof("Start mining at height %d with %d txs", lastBlock.Height, len(txs))

	var txList []*blockchain.Transaction
	for _, tx := range txs {
		txList = append(txList, &tx)
	}

	resultCh := make(chan *blockchain.Block, 1)

	go net.listenCompetingBlocks(ctxInternal, lastBlock.Height, resultCh)

	go func() {
		net.IsMining = true
		block, err := chain.MineBlock(
			txList,
			MinerAddress,
			net.HandleReoganizeTx,
			ctx,
		)

		if err != nil {
			log.Errorf("MineBlock error: %v", err)
			resultCh <- nil
			return
		}

		resultCh <- block
	}()

	select {
	case <-ctx.Done():
		log.Info("Mining cancelled by context")
		net.IsMining = false
		return
	case mined := <-resultCh:
		if mined == nil {
			log.Info("Mining stopped or failed.")
			for txID, tx := range MemoryPool.Queued {
				if !MemoryPool.HasPending(txID) {
					MemoryPool.Move(tx, memopool.MEMO_MOVE_FLAG_PENDING)
				}
			}
			ctx.Done()
			return
		}
		net.handleMinedBlock(mined)
		net.IsMining = false
		ctxInternal.Done()
	}
}

func (net *Network) listenCompetingBlocks(ctx context.Context, currentHeight int64, stop chan<- *blockchain.Block) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-net.competingBlockChan:
			if block.Height > currentHeight {
				log.Infof("Competing block detected (height %d), stop mining", block.Height)
				stop <- nil
				net.IsMining = false
				return
			}
		default:
			if !net.syncCompleted {
				stop <- nil
				return
			}
		}

	}
}

func (net *Network) handleMinedBlock(block *blockchain.Block) {
	utxo := blockchain.UTXOSet{Blockchain: net.Blockchain}
	utxo.Compute()

	log.Infof("New Block Mined: %x", block.Hash)

	net.Blocks <- block

	keys := slices.Collect(maps.Keys(MemoryPool.Queued))

	for _, key := range keys {
		MemoryPool.RemoveFromAll(key)
	}

}
