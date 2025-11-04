package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"core-blockchain/memopool"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (net *Network) HandleRequestSync() {
	listPeer := net.FullNodesChannel.ListPeers()

	for len(listPeer) == 0 {
		listPeer = net.FullNodesChannel.ListPeers()
	}

	log.Infof("Network sync: %d peer(s) connected — starting synchronization...", len(listPeer))

	time.Sleep(10 * time.Second)

	locator, err := net.Blockchain.GetBlockLocator()
	if err != nil {
		log.Error("Network sync aborted: failed to get block locator →", err)
		return
	}

	log.Infof("Network sync: successfully retrieved block locator with %d entries", len(locator))

	payload := NetHeaderLocator{
		SendFrom: net.Host.ID().String(),
		Locator:  locator,
	}

	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		if net.syncCompleted {
			break
		}

		log.Infof("Network sync: broadcasting header locator to %d peer(s)...", len(net.FullNodesChannel.ListPeers()))

		net.Gossip.Broadcast(
			net.FullNodesChannel.ListPeers(),
			[]string{
				net.Host.ID().String(),
			},
			func(p peer.ID) {
				log.Infof("Network sync: sending header locator to peer %s", p.String())

				net.SendHeaderLocator(
					p.String(),
					payload,
				)
			},
		)
	}

	log.Info("Network sync: header locator broadcast completed — waiting for peer responses...")
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
		log.Error("❌ Transaction handling aborted — failed to decode Tx payload")
		return
	}

	buf := bytes.NewBuffer(payload.Transaction)
	newTx := blockchain.DeserializeTxData(buf)

	if !net.syncCompleted {
		net.Gossip.Broadcast(
			net.FullNodesChannel.ListPeers(),
			[]string{
				net.Host.ID().String(),
			},
			func(p peer.ID) {
				net.SendTx(
					p.String(),
					newTx,
				)
				net.Gossip.MarkSeen(hex.EncodeToString(newTx.ID), p.String())
			},
		)

		return
	}

	txID := hex.EncodeToString(newTx.ID)
	if !MemoryPool.HasTX(txID) {
		txInfo := memopool.GetTxInfo(newTx, net.Blockchain)
		if txInfo == nil {
			log.Error("🚫 Transaction rejected — invalid or failed validation")
			return
		}

		MemoryPool.Add(*txInfo)
		log.Infof("✅ Transaction accepted and added to mempool — current size: %d", len(MemoryPool.Pending))

	}

	if !net.Gossip.HasSeen(txID, payload.SendFrom) {
		net.Gossip.Broadcast(
			net.FullNodesChannel.ListPeers(),
			[]string{
				net.Host.ID().String(),
				payload.SendFrom,
			},
			func(p peer.ID) {
				net.SendTx(p.String(), newTx)
				net.Gossip.MarkSeen(txID, payload.SendFrom)
			},
		)
		log.Infof("📡 Broadcasted transaction to peers (origin: %s)", payload.SendFrom)
	}
}

func (net *Network) HandleGetBlockDataSync(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetBlockSync

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("[Sync]❌ Aborted: cannot decode block data from peer %s", content.SendFrom)
		return
	}

	bestPeer := net.syncManager.GetTargetPeer()
	if payload.SendFrom != bestPeer.ID {
		log.Warnf("[Sync]⚠️ Ignored block data from %s (expected best peer %s)", payload.SendFrom, bestPeer.ID)
		return
	}

	log.Infof("[Sync]📥 Received %d blocks from best peer %s (peer height=%d)", len(payload.Blocks), payload.SendFrom, bestPeer.Height)

	isBadPeer := false

	for i := len(payload.Blocks) - 1; i >= 0; i-- {
		block := payload.Blocks[i]

		hasBlock, err := net.Blockchain.HasBlock(block.Hash)
		if err != nil {
			log.Errorf("[Sync]❌ Aborted: failed to check block %x at height %d: %v", block.Hash[:6], block.Height, err)
			return
		}

		if hasBlock {
			log.Infof("[Sync]⚠️ Skipped block %d (%x): already in chain", block.Height, block.Hash[:6])
			continue
		}

		if err := net.Blockchain.AddBlock(&block, net.HandleReoganizeTx); err != nil {
			log.Errorf("[Sync]❌ Aborted: cannot add block %d (%x): %v", block.Height, block.Hash[:6], err)
			isBadPeer = true
			break
		}

		log.Infof("[Sync] ✅ Added block %d (%x) to chain", block.Height, block.Hash[:6])
	}

	locator, err := net.Blockchain.GetBlockLocator()
	if err != nil {
		log.Errorf("[Sync]❌ Failed to build block locator: %v", err)
		return
	}

	excludePeerIDs := []string{
		net.Host.ID().String(),
	}

	if isBadPeer {
		net.syncManager.ClearTarget()
		excludePeerIDs = append(excludePeerIDs, payload.SendFrom)
	}

	time.Sleep(time.Second)

	net.Gossip.Broadcast(
		net.FullNodesChannel.ListPeers(),
		excludePeerIDs,
		func(p peer.ID) {
			net.SendHeaderLocator(p.String(), NetHeaderLocator{
				SendFrom: net.Host.ID().String(),
				Locator:  locator,
			})
		},
	)
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
		block, err := net.Blockchain.GetBlockMainChain(blockByte)
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
		bestPeer := net.syncManager.GetTargetPeer()

		if bestPeer != nil && bestPeer.Height >= bestHeight {
			log.Infof("✅ Header sync completed: peer %s chain matches local best height %d",
				payload.SendFrom, bestHeight)
			net.syncCompleted = true
		} else {
			net.syncManager.UpdatePeerStatus(
				payload.SendFrom,
				payload.BestHeight,
				big.NewInt(0),
			)

			bestPeer := net.syncManager.GetTargetPeer()

			if bestPeer != nil && bestPeer.Height >= bestHeight {
				log.Infof("✅ Header sync completed: peer %s chain matches local best height %d",
					payload.SendFrom, bestHeight)
				net.syncCompleted = true
			} else {
				log.Infof("Header sync: received empty headers from peer %s (local best height: %d, peer best height: %d)",
					payload.SendFrom, bestHeight, payload.BestHeight)
			}
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

		if bestHeight < bestPeer.Height {
			net.SendGetDataSync(payload.SendFrom, NetGetDataSync{
				SendFrom: net.Host.ID().String(),
				Hashes:   hashes,
			})
		} else {
			net.syncCompleted = true
			log.Infof("✅ Header sync completed: peer %s chain matches local best height %d",
				payload.SendFrom, bestHeight)
		}

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
		block, err := net.Blockchain.GetBlockMainChain(hashByte)
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
		log.Errorf("[HandleGetBlockData]❌ Failed to decode block data request from peer %s: %v", content.SendFrom, err)
		return
	}

	block := &blockchain.Block{}
	block = blockchain.DeserializeBlockData(payload.Block)

	log.Infof("[HandleGetBlockData]📥 Received block data %x (height=%d) from peer %s", block.Hash[:6], block.Height, payload.SendFrom)

	hasBlock, err := net.Blockchain.HasBlock(block.Hash)
	if err != nil {
		log.Errorf("[HandleGetBlockData]❌ Error checking block existence %x: %v", block.Hash[:6], err)
		return
	}

	if hasBlock {
		log.Warnf("[HandleGetBlockData]⚠️ Block %x (height=%d) already exists", block.Hash[:6], block.Height)
		if !net.Gossip.HasSeen(hex.EncodeToString(block.Hash), payload.SendFrom) {
			log.Infof("[HandleGetBlockData]🔄 Rebroadcasting header for existing block %x", block.Hash[:6])
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
		log.Errorf("[HandleGetBlockData]❌ Failed to add block %x: %v", block.Hash[:6], err)
		return
	}

	log.Infof("[HandleGetBlockData]✅ Block %d (%x) added to chain", block.Height, block.Hash[:6])

	for _, tx := range block.Transactions {
		txID := hex.EncodeToString(tx.ID)
		MemoryPool.RemoveFromAll(txID)
	}

	if net.Miner && net.IsMining {
		log.Infof("[HandleGetBlockData]⛏️ Competing block received while mining: %x", block.Hash[:6])
		net.competingBlockChan <- block
	}

	log.Infof("[HandleGetBlockData]Broadcasting block header hash=%x height=%d to peers (excluding sender=%s)", block.Hash[:6], block.Height, payload.SendFrom)
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
	log.Debugf("[HandleGetBlockData]Completed handling GetBlockData from peer=%s for block hash=%x height=%d", content.SendFrom, block.Hash[:6], block.Height)
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

	block, err := net.Blockchain.GetBlockMainChain(payload.Hash)
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
		log.Error("Failed to decode header request payload")
		return
	}

	exists, err := net.Blockchain.HasBlock(payload.Hash)
	if err != nil {
		log.Errorf("Error checking block existence: %v", err)
		return
	}

	lastBlock, err := net.Blockchain.GetLastBlock()
	if err != nil {
		log.Errorf("Error retrieving last block: %v", err)
		return
	}

	shortHash := fmt.Sprintf("%x", payload.Hash[:6])
	peerID := payload.SendFrom

	if exists {
		log.Infof("Received header for existing block [%s] from peer [%s]", shortHash, peerID)

		blockHash := hex.EncodeToString(payload.Hash)
		if !net.Gossip.HasSeen(blockHash, peerID) {
			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{
					peerID,
					net.Host.ID().String(),
				},
				func(p peer.ID) {
					net.SendHeader(p.String(), &NetHeader{
						Hash:     payload.Hash,
						Height:   payload.Height,
						PrevHash: payload.PrevHash,
						SendFrom: net.Host.ID().String(),
					})
					net.Gossip.MarkSeen(blockHash, peerID)
				},
			)
			log.Debugf("Broadcasted header [%s] to peers (origin: %s)", shortHash, peerID)
		} else {
			log.Debugf("Header [%s] from peer [%s] already seen — ignoring", shortHash, peerID)
		}

	} else {
		if bytes.Equal(lastBlock.Hash, payload.PrevHash) {
			net.syncCompleted = true
			log.Infof("Header [%s] links to last known block — requesting full block from peer [%s]", shortHash, peerID)
			net.SendGetData(peerID, NetGetData{
				SendFrom: net.Host.ID().String(),
				Height:   payload.Height,
				Hash:     payload.Hash,
			})
		} else if payload.Height > lastBlock.Height+1 {
			log.Warnf("Peer [%s] is ahead — header [%s] indicates height %d vs local %d. Requesting block locator...",
				peerID, shortHash, payload.Height, lastBlock.Height)

			locator, err := net.Blockchain.GetBlockLocator()
			if err != nil {
				log.Errorf("Error building block locator: %v", err)
				return
			}

			net.syncCompleted = false
			net.SendHeaderLocator(peerID, NetHeaderLocator{
				SendFrom: net.Host.ID().String(),
				Locator:  locator,
			})
			log.Infof("Sent header locator request to peer [%s]", peerID)
		} else {
			log.Infof("Received non-continuous header [%s] from peer [%s] — no action taken", shortHash, peerID)
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
		log.Errorf("TxPool sync aborted: failed to decode full transaction data from peer: %v", err)
		return
	}

	if len(payload.Transactions) <= 0 {
		return
	}

	for i, tx := range payload.Transactions {
		txID := hex.EncodeToString(tx.ID)
		txInfo := memopool.GetTxInfo(&tx, net.Blockchain)
		if txInfo != nil && !MemoryPool.HasTX(txID) {
			MemoryPool.Add(*txInfo)
			log.Debugf("TxPool sync: [%d/%d] added transaction %s to local mempool",
				i+1, len(payload.Transactions), txID)
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
		log.Errorf("TxPool sync aborted: failed to decode get transaction request from peer: %v", err)
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
		return
	}

	txHashes := make([][]byte, 0)

	for _, txHash := range payload.TxHashes {
		txID := hex.EncodeToString(txHash)
		if !MemoryPool.HasTX(txID) {
			txHashes = append(txHashes, txHash)
		}
	}

	if len(txHashes) > 0 {

		net.SendGetDataTransaction(
			payload.SendFrom,
			NetGetDataTransaction{
				SendFrom: net.Host.ID().String(),
				TxHashes: txHashes,
			},
		)

	}

}

func (net *Network) HandleGetTxFromPool(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload TxFromPool

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("❌ Failed to decode 'GetTxFromPool' request: %v", err)
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
				log.Infof("Transaction %s and fee %f added to mempool, broadcasting to peers...", txHash, txInfo.Fee)

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
