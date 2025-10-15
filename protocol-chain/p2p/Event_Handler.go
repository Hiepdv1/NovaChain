package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"core-blockchain/memopool"
	"encoding/gob"
	"encoding/hex"
	"math/big"

	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (net *Network) HandleRequestSync() {
	listPeer := net.FullNodesChannel.ListPeers()
	if len(listPeer) == 0 {
		return
	}

	log.Info("Starting network synchronization...")

	net.syncCompleted = false

	locator, err := net.Blockchain.GetBlockLocator()
	if err != nil {
		log.Error("Failed to get block locator: ", err)
		return
	}

	payload := NetHeaderLocator{
		SendFrom: net.Host.ID().String(),
		Locator:  locator,
	}

	net.Gossip.Broadcast(
		net.FullNodesChannel.ListPeers(),
		[]string{
			net.Host.ID().String(),
		},
		func(p peer.ID) {
			log.Infof("Sending header locator to peer: %s", p.String())

			net.SendHeaderLocator(
				p.String(),
				payload,
			)
		},
	)

}

func (net *Network) HandleReoganizeTx(txs []*blockchain.Transaction) {
	net.Gossip.Broadcast(
		net.FullNodesChannel.ListPeers(),
		[]string{
			net.Host.ID().String(),
		},
		func(p peer.ID) {
			for _, tx := range txs {
				net.SendTx(p.String(), tx)
			}
		},
	)
}

func (net *Network) HandleTxMining(content *ChannelContent) {
	buff := new(bytes.Buffer)
	var payload NetTxMining

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding transaction request: ", err)
		return
	}
	for _, tx := range payload.Txs {
		txInfo := memopool.GetTxInfo(&tx, net.Blockchain)
		if txInfo == nil {
			continue
		}
		MemoryPool.Move(*txInfo, memopool.MEMO_MOVE_FLAG_QUEUED)
	}
}

func (net *Network) HandleTx(content *ChannelContent) {
	buff := new(bytes.Buffer)
	var payload NetTx

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buff)
	if err := dec.Decode(&payload); err != nil {
		log.Error("Tx sync aborted: failed to decode request")
		return
	}

	buf := bytes.NewBuffer(payload.Transaction)

	newTx := blockchain.DeserializeTxData(buf)

	txInfo := memopool.GetTxInfo(newTx, net.Blockchain)
	if txInfo == nil {
		log.Error("Tx rejected: invalid transaction data")
		return
	}

	MemoryPool.Add(*txInfo)
	log.Infof("Tx accepted → Mempool size: %d", len(MemoryPool.Pending))

	if !net.Gossip.HasSeen(hex.EncodeToString(newTx.ID), payload.SendFrom) {
		net.Gossip.Broadcast(
			net.FullNodesChannel.ListPeers(),
			[]string{
				net.Host.ID().String(),
				payload.SendFrom,
			},
			func(p peer.ID) {
				net.SendTx(p.String(), newTx)
			},
		)
		log.Info("Tx broadcasted to peers")
	}
}

func (net *Network) HandleGetBlockDataSync(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetBlockSync

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("[Sync] ❌ Aborted: cannot decode block data from peer %s", content.SendFrom)
		return
	}

	bestPeer := net.syncManager.GetTargetPeer()
	if payload.SendFrom != bestPeer.ID {
		log.Warnf("[Sync] ⚠️ Ignored block data from %s (expected best peer %s)", payload.SendFrom, bestPeer.ID)
		return
	}

	log.Infof("[Sync] 📥 Received %d blocks from best peer %s (peer height=%d)", len(payload.Blocks), payload.SendFrom, bestPeer.Height)

	for i := len(payload.Blocks) - 1; i >= 0; i-- {
		block := payload.Blocks[i]

		hasBlock, err := net.Blockchain.HasBlock(block.Hash)
		if err != nil {
			log.Errorf("[Sync] ❌ Aborted: failed to check block %x at height %d: %v", block.Hash[:6], block.Height, err)
			return
		}

		if hasBlock {
			log.Infof("[Sync] ⚠️ Skipped block %d (%x): already in chain", block.Height, block.Hash[:6])
			continue
		}

		if err := net.Blockchain.AddBlock(&block, net.HandleReoganizeTx); err != nil {
			log.Errorf("[Sync] ❌ Aborted: cannot add block %d (%x): %v", block.Height, block.Hash[:6], err)
			return
		}

		log.Infof("[Sync] ✅ Added block %d (%x) to chain", block.Height, block.Hash[:6])
	}

	UTXO := blockchain.UTXOSet{Blockchain: net.Blockchain}
	if err := UTXO.Compute(); err != nil {
		log.Errorf("[Sync] ❌ Failed to recompute UTXO: %v", err)
		return
	}
	log.Info("[Sync] 🗂️ UTXO recomputed after block batch")

	currentHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("[Sync] ❌ Failed to get best height: %v", err)
		return
	}

	if currentHeight < bestPeer.Height {
		locator, err := net.Blockchain.GetBlockLocator()
		if err != nil {
			log.Errorf("[Sync] ❌ Failed to build block locator: %v", err)
			return
		}
		net.SendHeaderLocator(payload.SendFrom, NetHeaderLocator{
			SendFrom: net.Host.ID().String(),
			Locator:  locator,
		})
		net.syncCompleted = false
		log.Infof("[Sync] ⏩ Local height=%d < Peer height=%d → requesting more blocks...", currentHeight, bestPeer.Height)
	} else {
		net.syncCompleted = true
		log.Infof("[Sync] 🎉 Node fully synced with peer %s at height=%d", payload.SendFrom, currentHeight)
	}
}

func (net *Network) HandleGetDataSync(content *ChannelContent) {
	buf := new(bytes.Buffer)

	buf.Write(content.Payload[commandLength:])
	var payload NetGetDataSync

	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Error("Sync aborted: failed to decode request")
		return
	}

	blocks := make([]blockchain.Block, 0)
	for _, blockByte := range payload.Hashes {
		block, err := net.Blockchain.GetBlock(blockByte)
		if err != nil {
			log.Error("Failed to retrieve block for sync")
			return
		}
		blocks = append(blocks, BlockForNetwork(block))
	}

	bestHeight, err := net.Blockchain.GetBestHeight()

	if err != nil {
		log.Errorf("Get Best Height with error: %v", err)
		return
	}

	net.SendBlockDataSync(payload.SendFrom, NetBlockSync{
		SendFrom:   net.Host.ID().String(),
		BestHeight: bestHeight,
		Blocks:     blocks,
	})

	log.Info("Sent requested blocks to peer")
}

func (net *Network) HandleGetHeaderSync(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetHeaders

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Error("Header sync aborted: failed to decode headers from peer")
		return
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("Get Best Height with error: %v", err)
		return
	}

	if len(payload.Data) < 1 {
		net.syncManager.UpdatePeerStatus(
			payload.SendFrom,
			payload.BestHeight,
			big.NewInt(0),
		)
		bestPeer := net.syncManager.GetTargetPeer()

		if bestPeer != nil && bestPeer.Height == bestHeight {
			log.Infof("✅ Header sync completed: peer %s chain matches local best height %d",
				payload.SendFrom, bestHeight)
			net.syncCompleted = true
		} else {
			log.Infof("Header sync: received empty headers from peer %s (local best height: %d, peer best height: %d)",
				payload.SendFrom, bestHeight, payload.BestHeight)
		}

		return
	}

	block := payload.Data[len(payload.Data)-1]

	log.Infof("Received %d blocks", len(payload.Data))

	exist, err := net.Blockchain.HasBlock(block.PrevHash)
	if err != nil {
		log.Error("Header sync failed: unable to check previous hash")
		return
	}

	if !exist {
		log.Warnf("Header sync skipped: previous hash %x not found in local chain - height: %d", block.PrevHash, block.Height)
		return
	}

	peerTotalWork := new(big.Int)
	hashes := make([][]byte, 0)

	for _, header := range payload.Data {
		work := net.Blockchain.CalcWork(header.Nbits)
		peerTotalWork = new(big.Int).Add(peerTotalWork, work)
		hashes = append(hashes, header.Hash)
	}

	net.syncManager.UpdatePeerStatus(
		payload.SendFrom,
		payload.BestHeight,
		peerTotalWork,
	)

	bestPeer := net.syncManager.GetTargetPeer()

	if bestPeer.ID == payload.SendFrom {
		log.Infof("Header sync: received %d headers from best peer %s, requesting full blocks for these headers",
			len(hashes), payload.SendFrom)

		net.SendGetDataSync(payload.SendFrom, NetGetDataSync{
			SendFrom: net.Host.ID().String(),
			Hashes:   hashes,
		})

	} else {
		log.Infof("Header sync skipped: peer %s chain has less work than current best peer",
			payload.SendFrom)
	}
}

func (net *Network) HandleGetHeaderLocator(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetHeaderLocator

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Error("Header sync aborted: failed to decode locator request")
		return
	}

	var commonBlock *blockchain.Block
	for _, hashByte := range payload.Locator {
		block, err := net.Blockchain.GetBlock(hashByte)
		if err != nil {
			continue
		}
		commonBlock = &block
		break
	}

	if commonBlock == nil {
		log.Error("Header sync skipped: no common block found")
		return
	}

	log.Infof("Block common with height %d", commonBlock.Height)

	blocks, err := net.Blockchain.GetBlockRange(commonBlock.Hash, MAX_HEADERS_PER_MSG)
	if err != nil {
		log.Error("Header sync failed: cannot fetch block range")
		return
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("Get Best Height with error: %v", err)
		return
	}

	if len(blocks) == 0 {
		net.SendHeaders(payload.SendFrom, NetHeaders{
			SendFrom:   net.Host.ID().String(),
			BestHeight: bestHeight,
			Data:       []NetHeadersData{},
		})
		return
	}

	data := make([]NetHeadersData, 0)
	for _, block := range blocks {
		data = append(data, NetHeadersData{
			Height:   block.Height,
			PrevHash: block.PrevHash,
			Hash:     block.Hash,
			Nbits:    block.NBits,
		})
	}

	net.SendHeaders(payload.SendFrom, NetHeaders{
		SendFrom:   net.Host.ID().String(),
		BestHeight: bestHeight,
		Data:       data,
	})

	log.Info("Sent headers to peer for synchronization")
}

func (net *Network) HandleGetBlockData(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetBlock

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("[HandleGetBlockData] ❌ Failed to decode block data request from peer %s: %v", content.SendFrom, err)
		return
	}

	block := &blockchain.Block{}
	block = blockchain.DeserializeBlockData(payload.Block)

	log.Infof("[HandleGetBlockData] 📥 Received block data %x (height=%d) from peer %s", block.Hash[:6], block.Height, payload.SendFrom)

	hasBlock, err := net.Blockchain.HasBlock(block.Hash)
	if err != nil {
		log.Errorf("[HandleGetBlockData] ❌ Error checking block existence %x: %v", block.Hash[:6], err)
		return
	}

	if hasBlock {
		log.Warnf("[HandleGetBlockData] ⚠️ Block %x (height=%d) already exists", block.Hash[:6], block.Height)
		if !net.Gossip.HasSeen(hex.EncodeToString(block.Hash), payload.SendFrom) {
			log.Infof("[HandleGetBlockData] 🔄 Rebroadcasting header for existing block %x", block.Hash[:6])
			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{net.Host.ID().String(), payload.SendFrom},
				func(p peer.ID) {
					net.SendHeader(p.String(), &NetHeader{
						Hash:     block.Hash,
						Height:   block.Height,
						PrevHash: block.PrevHash,
						SendFrom: net.Host.ID().String(),
					})
					net.Gossip.MarkSeen(hex.EncodeToString(block.Hash), p.String())
				},
			)
		}
		return
	}

	err = net.Blockchain.AddBlock(block, net.HandleReoganizeTx)
	if err != nil {
		log.Errorf("[HandleGetBlockData] ❌ Failed to add block %x: %v", block.Hash[:6], err)
		return
	}

	log.Infof("[HandleGetBlockData] ✅ Block %d (%x) added to chain", block.Height, block.Hash[:6])

	for _, tx := range block.Transactions {
		txID := hex.EncodeToString(tx.ID)
		MemoryPool.RemoveFromAll(txID)
	}

	if net.Miner && net.IsMining {
		log.Infof("[HandleGetBlockData] ⛏️ Competing block received while mining: %x", block.Hash[:6])
		net.competingBlockChan <- block
	}

	UTXO := blockchain.UTXOSet{Blockchain: net.Blockchain}
	err = UTXO.Compute()
	if err != nil {
		log.Errorf("[HandleGetBlockData] ❌ Failed to recompute UTXO after block %x: %v", block.Hash[:6], err)
		return
	}

	log.Infof("[HandleGetBlockData] 🗂️ UTXO set recomputed successfully after block %x", block.Hash[:6])

	log.Infof("[HandleGetBlockData] Broadcasting block header hash=%x height=%d to peers (excluding sender=%s)", block.Hash[:6], block.Height, payload.SendFrom)
	net.Gossip.Broadcast(
		net.FullNodesChannel.ListPeers(),
		[]string{net.Host.ID().String(), payload.SendFrom},
		func(p peer.ID) {
			net.SendHeader(p.String(), &NetHeader{
				Hash:     block.Hash,
				Height:   block.Height,
				PrevHash: block.PrevHash,
				SendFrom: net.Host.ID().String(),
			})
			net.Gossip.MarkSeen(hex.EncodeToString(block.Hash), p.String())
		},
	)
	log.Debugf("[HandleGetBlockData] Completed handling GetBlockData from peer=%s for block hash=%x height=%d", content.SendFrom, block.Hash[:6], block.Height)
}

func (net *Network) HandleGetData(content *ChannelContent) {

	buf := new(bytes.Buffer)
	var payload NetGetData

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Failed to decode GetData request")
		return
	}

	block, err := net.Blockchain.GetBlock(payload.Hash)
	if err != nil {
		log.Errorf("Block %x not found in chain", payload.Hash[:6])
		return
	}

	block = BlockForNetwork(block)
	net.SendFullBlock(payload.SendFrom, &block)

	log.Infof("📦 Sent block %d (%x) to peer %s", block.Height, block.Hash[:6], payload.SendFrom)
}

func (net *Network) HandleGetHeader(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetHeader

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Failed to decode header request")
		return
	}

	exists, err := net.Blockchain.HasBlock(payload.Hash)
	if err != nil {
		log.Error(err)
		return
	}

	lastBlock, err := net.Blockchain.GetLastBlock()
	if err != nil {
		log.Error(err)
		return
	}

	if exists {
		blockHash := hex.EncodeToString(payload.Hash)
		log.Infof("Header received for existing block %x", payload.Hash[:6])

		if !net.Gossip.HasSeen(blockHash, payload.SendFrom) {
			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{
					payload.SendFrom,
					net.Host.ID().String(),
				},
				func(p peer.ID) {
					net.SendHeader(p.String(), &NetHeader{
						Hash:     payload.Hash,
						Height:   payload.Height,
						PrevHash: payload.PrevHash,
						SendFrom: net.Host.ID().String(),
					})
					net.Gossip.MarkSeen(blockHash, payload.SendFrom)
				},
			)
		}

	} else {
		if bytes.Equal(lastBlock.Hash, payload.PrevHash) {
			log.Infof("Requesting full block %s from peer %s", hex.EncodeToString(payload.Hash), payload.SendFrom)
			net.SendGetData(payload.SendFrom, NetGetData{
				SendFrom: net.Host.ID().String(),
				Height:   payload.Height,
				Hash:     payload.Hash,
			})
		} else if payload.Height > lastBlock.Height+1 {
			log.Warnf("Header %x indicates peer is ahead (height %d vs local %d), requesting locator...", payload.Hash[:6], payload.Height, lastBlock.Height)
			locator, err := net.Blockchain.GetBlockLocator()
			if err != nil {
				log.Error(err)
				return
			}
			net.syncCompleted = false
			net.SendHeaderLocator(payload.SendFrom, NetHeaderLocator{
				SendFrom: net.Host.ID().String(),
				Locator:  locator,
			})
		}
	}
}

func (net *Network) HandleGetTransactions(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetTransactionData

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get transaction from pool request: ", err)
		return
	}

	if len(payload.Transactions) <= 0 {
		return
	}

	for _, tx := range payload.Transactions {
		txInfo := memopool.GetTxInfo(&tx, net.Blockchain)
		if txInfo != nil {
			MemoryPool.Add(*txInfo)
		}
	}
}

func (net *Network) HandleGetDataTx(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetGetDataTransaction

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get transaction from pool request: ", err)
		return
	}

	var txs []blockchain.Transaction

	for _, txIdBytes := range payload.TxHashes {
		txID := hex.EncodeToString(txIdBytes)
		tx := MemoryPool.GetTxByID(txID)
		if tx != nil {
			txs = append(txs, *tx)
		}
	}

	net.SendTransactions(
		payload.SendFrom,
		NetTransactionData{
			SendFrom:     net.Host.ID().String(),
			Transactions: txs,
		},
	)
}

func (net *Network) HandleGetTxPoolInv(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetInventoryTxs

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Error("TxPool sync aborted: failed to decode transaction inventory from peer")
		return
	}

	log.Infof("TxPool sync: received tx inventory with %d hashes from peer %s",
		len(payload.TxHashes), payload.SendFrom)

	txHashes := make([][]byte, 0)

	for _, txHash := range payload.TxHashes {
		txID := hex.EncodeToString(txHash)
		if !MemoryPool.HashTX(txID) {
			txHashes = append(txHashes, txHash)
			log.Debugf("TxPool sync: missing transaction hash %s, will request full tx from peer %s",
				txID, payload.SendFrom)
		} else {
			log.Tracef("TxPool sync: transaction hash %s already exists in local mempool, skip request", txID)
		}
	}

	if len(txHashes) > 0 {
		log.Infof("TxPool sync: requesting %d full transactions from peer %s",
			len(txHashes), payload.SendFrom)

		net.SendGetDataTransaction(
			payload.SendFrom,
			NetGetDataTransaction{
				SendFrom: net.Host.ID().String(),
				TxHashes: txHashes,
			},
		)
	} else {
		log.Infof("TxPool sync: all %d transactions already in mempool, nothing to request from peer %s",
			len(payload.TxHashes), payload.SendFrom)
	}
}

func (net *Network) HandleGetTxFromPool(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload *TxFromPool

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get transaction from pool request: ", err)
		return
	}

	if int64(len(MemoryPool.Pending)) >= payload.Count {
		txs := MemoryPool.GetTransactionHashes()
		net.SendTxPoolInv(payload.SendFrom, txs)
	} else {
		net.SendTxPoolInv(payload.SendFrom, [][]byte{})
	}
}
func HandleEvents(net *Network) {
	for {
		select {
		case block := <-net.Blocks:
			bHash := hex.EncodeToString(block.Hash)
			log.Infof("HandleEvents: received new block hash=%s height=%d, broadcasting header to peers...", bHash, block.Height)

			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{
					net.Host.ID().String(),
				},
				func(p peer.ID) {
					log.Infof("Broadcasting block header %s to peer %s", bHash, p.String())

					net.SendHeader(p.String(), &NetHeader{
						SendFrom: net.Host.ID().String(),
						Hash:     block.Hash,
						Height:   block.Height,
						PrevHash: block.PrevHash,
					})

					net.Gossip.MarkSeen(bHash, p.String())
				},
			)

		case txs := <-net.Transactions:
			log.Infof("HandleEvents: received %d new transactions, processing...", len(txs))

			for _, tx := range txs {
				txHash := hex.EncodeToString(tx.ID)
				txInfo := memopool.GetTxInfo(tx, net.Blockchain)
				if txInfo == nil {
					log.Errorf("Failed to get transaction info for tx=%s", txHash)
					continue
				}

				MemoryPool.Add(*txInfo)
				log.Infof("Transaction %s added to mempool, broadcasting to peers...", txHash)

				net.Gossip.Broadcast(
					net.FullNodesChannel.ListPeers(),
					[]string{
						net.Host.ID().String(),
					},
					func(p peer.ID) {
						log.Infof("Broadcasting tx %s to peer %s", txHash, p.String())
						net.SendTx(p.String(), tx)
						net.Gossip.MarkSeen(txHash, p.String())
					},
				)
			}
		}
	}
}
