package p2p

import (
	"context"
	blockchain "core-blockchain/core"
	"maps"
	"slices"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (net *Network) MinersEventLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	time.Sleep(10 * time.Second)

	go net.requestTransaction(ctx)
	net.miningLoop(ctx)
}

func (net *Network) requestTransaction(ctx context.Context) {
	log.Info("TxPool sync: starting periodic transaction request routine (interval = 10s)")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("TxPool sync: stopped periodic transaction request routine")
			return
		case <-ticker.C:
			peers := net.FullNodesChannel.ListPeers()
			if len(peers) == 0 {
				continue
			}

			net.Gossip.Broadcast(
				peers,
				[]string{net.Host.ID().String()},
				func(p peer.ID) {
					net.SendRequestTxFromPool(p.String(), TxFromPool{
						SendFrom: net.Host.ID().String(),
						Count:    1,
					})
				},
			)
		}
	}
}

func (net *Network) miningLoop(ctx context.Context) {
	net.IsMining = true
	for {
		select {
		case <-ctx.Done():
			log.Info("Mining loop stopped.")
			net.IsMining = false
			return
		default:
			if !net.syncCompleted {
				time.Sleep(3 * time.Second)
				continue
			}

			txs := MemoryPool.SelectHighFeeTx()
			if len(txs) == 0 {
				txs = map[string]blockchain.Transaction{}
			}

			miningCtx, cancel := context.WithCancel(ctx)
			net.tryMine(miningCtx, cancel, txs)

			time.Sleep(time.Second)
		}
	}
}

func (net *Network) tryMine(ctx context.Context, cancel context.CancelFunc, txs map[string]blockchain.Transaction) {
	if !net.syncCompleted {
		return
	}

	log.Info("Starting mining new block...")
	net.mineBlock(ctx, cancel, txs)
}

func (net *Network) mineBlock(ctx context.Context, cancel context.CancelFunc, txs map[string]blockchain.Transaction) {
	chain := net.Blockchain
	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		log.Errorf("Get last block error: %v", err)
		return
	}

	log.Infof("Mining at height %d with %d txs", lastBlock.Height+1, len(txs))
	txList := make([]*blockchain.Transaction, 0)
	for _, tx := range txs {
		txList = append(txList, &tx)
	}

	resultCh := make(chan *blockchain.Block, 1)

	go net.listenCompetingBlocks(ctx, lastBlock.Height, cancel)

	go net.startMiningWorker(ctx, resultCh, txList)

	select {
	case <-ctx.Done():
		log.Warn("Mining canceled.")
	case mined := <-resultCh:
		if mined != nil {
			net.handleMinedBlock(mined)
		}
		cancel()
	}
}

func (net *Network) startMiningWorker(ctx context.Context, resultCh chan<- *blockchain.Block, txList []*blockchain.Transaction) {
	chain := net.Blockchain

	log.Infof("[Miner] start mining...")
	block, err := chain.MineBlock(txList, MinerAddress, net.HandleReoganizeTx, ctx)
	if err != nil {
		log.Warnf("[Miner] MineBlock: %v", err)
		resultCh <- nil
		return
	}

	resultCh <- block
	log.Infof("[Miner] found valid block at height %d!", block.Height)

}

func (net *Network) listenCompetingBlocks(ctx context.Context, currentHeight int64, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-net.competingBlockChan:
			if block.Height >= currentHeight {
				log.Warnf("Competing block detected (height %d), stop mining", block.Height)
				cancel()
				return
			}
		}
	}
}

func (net *Network) handleMinedBlock(block *blockchain.Block) {
	log.Infof("New Block Mined: %x", block.Hash)
	net.Blocks <- block

	keys := slices.Collect(maps.Keys(MemoryPool.Queued))
	for _, key := range keys {
		MemoryPool.RemoveFromAll(key)
	}
}
