package p2p

import (
	"bytes"
	"context"
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
	const logName = "[SYNC::REQUEST]"

	listPeer := net.FullNodesChannel.ListPeers()

	if len(listPeer) == 0 {
		return
	}

	log.Infof("%s %d peer(s) connected — starting synchronization...", logName, len(listPeer))

	time.Sleep(10 * time.Second)

	locator, err := net.Blockchain.GetBlockLocator()
	if err != nil {
		log.Errorf("%s Aborted: failed to get block locator → %v", logName, err)
		return
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("%s Aborted: failed to get best height → %v", logName, err)
		return
	}

	log.Infof("%s Successfully retrieved block locator with %d entries", logName, len(locator))

	payload := NetHeaderLocator{
		SendFrom:   net.Host.ID().String(),
		BestHeight: bestHeight,
		Locator:    locator,
	}

	net.Gossip.Broadcast(
		listPeer,
		[]string{net.Host.ID().String()},
		func(p peer.ID) {
			log.Infof("%s Sending header locator to peer %s", logName, p.String())
			net.SendHeaderLocator(p.String(), payload)
		},
	)

	log.Infof("%s Header locator broadcast completed — waiting for peer responses...", logName)
}

func (net *Network) HandleRequestGossipPeer(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			net.Gossip.Broadcast(
				net.FullNodesChannel.ListPeers(),
				[]string{net.Host.ID().String()},
				func(p peer.ID) {

					net.SendRequestGossipPeer(
						p.String(),
						NetRequestGossipPeer{
							SendFrom: net.Host.ID().String(),
							Count:    5,
						},
					)
				},
			)
		case <-ctx.Done():
			return
		}
	}
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
	const logName = "[SYNC::BLOCK_DATA]"

	buf := new(bytes.Buffer)
	var payload NetBlockSync

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("%s Aborted: cannot decode block data from peer %s", logName, content.SendFrom)
		return
	}

	bestPeer := net.syncManager.GetTargetPeer()
	if payload.SendFrom != bestPeer.ID {
		log.Warnf("%s Ignored block data from %s (expected best peer %s)", logName, payload.SendFrom, bestPeer.ID)
		return
	}

	log.Infof("%s Received %d blocks from best peer %s (peer height=%d)", logName, len(payload.Blocks), payload.SendFrom, bestPeer.Height)

	isBadPeer := false

	for i := len(payload.Blocks) - 1; i >= 0; i-- {
		block := payload.Blocks[i]

		if bestPeer.ID != payload.SendFrom {
			log.Warnf("%s Peer %s no longer matches best peer %s, aborting loop", logName, payload.SendFrom, bestPeer.ID)
			break
		}

		hasBlock, err := net.Blockchain.HasBlock(block.Hash)
		if err != nil {
			log.Errorf("%s Failed to check block %x at height %d: %v", logName, block.Hash[:6], block.Height, err)
			return
		}

		if hasBlock {
			log.Infof("%s Skipped block %d (%x): already exists", logName, block.Height, block.Hash[:6])
			continue
		}

		if err := net.Blockchain.AddBlock(&block, net.HandleReoganizeTx); err != nil {
			log.Errorf("%s Cannot add block %d (%x): %v", logName, block.Height, block.Hash[:6], err)
			isBadPeer = true
			break
		}

		log.Infof("%s Added block %d (%x) to chain", logName, block.Height, block.Hash[:6])
	}

	locator, err := net.Blockchain.GetBlockLocator()
	if err != nil {
		log.Errorf("%s Failed to build block locator: %v", logName, err)
		return
	}

	excludePeerIDs := []string{net.Host.ID().String()}

	if isBadPeer {
		net.syncManager.ClearTarget()
		excludePeerIDs = append(excludePeerIDs, payload.SendFrom)
		log.Warnf("%s Marked peer %s as bad and cleared target", logName, payload.SendFrom)
	}

	time.Sleep(time.Second)

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("%s Failed to get best height: %v", logName, err)
		return
	}

	log.Infof("%s Broadcasting header locator (best height=%d)", logName, bestHeight)
	net.Gossip.Broadcast(
		net.FullNodesChannel.ListPeers(),
		excludePeerIDs,
		func(p peer.ID) {
			net.SendHeaderLocator(p.String(), NetHeaderLocator{
				SendFrom:   net.Host.ID().String(),
				BestHeight: bestHeight,
				Locator:    locator,
			})
			log.Debugf("%s Sent header locator to peer %s", logName, p.String())
		},
	)
}

func (net *Network) HandleGetDataSync(content *ChannelContent) {
	const logName = "[SYNC::GET_DATA]"

	buf := new(bytes.Buffer)
	buf.Write(content.Payload[commandLength:])
	var payload NetGetDataSync

	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("%s Aborted: failed to decode data request from peer %s", logName, content.SendFrom)
		return
	}

	log.Infof("%s Received data request from peer %s (%d block hashes)", logName, payload.SendFrom, len(payload.Hashes))

	blocks := make([]blockchain.Block, 0)
	for _, blockByte := range payload.Hashes {
		block, err := net.Blockchain.GetBlockMainChain(blockByte)
		if err != nil {
			log.Errorf("%s Failed to retrieve block for hash %x: %v", logName, blockByte[:6], err)
			return
		}
		blocks = append(blocks, BlockForNetwork(block))
		log.Debugf("%s Prepared block %d (%x) for response", logName, block.Height, block.Hash[:6])
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("%s Failed to get best height: %v", logName, err)
		return
	}

	net.SendBlockDataSync(payload.SendFrom, NetBlockSync{
		SendFrom:   net.Host.ID().String(),
		BestHeight: bestHeight,
		Blocks:     blocks,
	})

	log.Infof("%s Sent %d requested blocks to peer %s (best height=%d)", logName, len(blocks), payload.SendFrom, bestHeight)
}

func (net *Network) HandleGetHeaderSync(content *ChannelContent) {
	const logName = "[SYNC::HEADER_SYNC]"

	buf := new(bytes.Buffer)
	var payload NetHeaders

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("%s Decode failed: invalid header data from peer", logName)
		return
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("%s Failed to get best height: %v", logName, err)
		return
	}

	if len(payload.Data) < 1 {
		bestPeer := net.syncManager.GetTargetPeer()

		if bestPeer != nil && bestPeer.Height >= bestHeight {
			log.Infof("%s Completed: peer %s chain matches local best height (%d)",
				logName, payload.SendFrom, bestHeight)
			net.syncCompleted = true
		} else {
			net.syncManager.UpdatePeerStatus(payload.SendFrom, payload.BestHeight, big.NewInt(0))
			bestPeer = net.syncManager.GetTargetPeer()

			if bestPeer != nil && bestPeer.Height >= bestHeight {
				log.Infof("%s Completed: peer %s at height %d (local best: %d)",
					logName, payload.SendFrom, bestPeer.Height, bestHeight)
				net.syncCompleted = true
			} else if bestHeight >= bestPeer.Height {
				log.Infof("%s Completed: local best height %d matches peer %s",
					logName, bestHeight, payload.SendFrom)
				net.syncCompleted = true
			} else {
				log.Infof("%s No new headers from peer %s (local best: %d, peer best: %d)",
					logName, payload.SendFrom, bestHeight, payload.BestHeight)
			}
		}
		return
	}

	block := payload.Data[len(payload.Data)-1]
	log.Infof("%s Received %d headers from peer %s", logName, len(payload.Data), payload.SendFrom)

	chainBlock, err := net.Blockchain.GetBlock(block.PrevHash)
	if err != nil {
		log.Errorf("%s Failed to validate previous hash from received headers", logName)
		return
	}

	peerTotalWork := new(big.Int).Set(chainBlock.NChainWork)
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
		log.Infof("%s Processing %d headers from best peer %s, requesting corresponding blocks",
			logName, len(hashes), payload.SendFrom)

		if bestHeight < bestPeer.Height {
			net.SendGetDataSync(payload.SendFrom, NetGetDataSync{
				SendFrom: net.Host.ID().String(),
				Hashes:   hashes,
			})
		} else {
			net.syncCompleted = true
			log.Infof("%s Completed: peer %s chain matches local best height (%d)",
				logName, payload.SendFrom, bestHeight)
		}
	} else {
		log.Infof("%s Skipped: peer %s chain has less work than best peer", logName, payload.SendFrom)
	}
}

func (net *Network) HandleGetHeaderLocator(content *ChannelContent) {
	const logName = "[SYNC::HEADER_LOCATOR]"

	buf := new(bytes.Buffer)
	var payload NetHeaderLocator

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&payload); err != nil {
		log.Errorf("%s Decode failed: invalid locator request from peer", logName)
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
		log.Warnf("%s No common ancestor block found with requesting peer", logName)
		return
	}

	log.Infof("%s Found common block at height %d", logName, commonBlock.Height)

	blocks, err := net.Blockchain.GetBlockRange(commonBlock.Hash, MAX_HEADERS_PER_MSG)
	if err != nil {
		log.Errorf("%s Failed to fetch block range for headers", logName)
		return
	}

	bestHeight, err := net.Blockchain.GetBestHeight()
	if err != nil {
		log.Errorf("%s Failed to get best height: %v", logName, err)
		return
	}

	if payload.BestHeight > bestHeight {
		locator, err := net.Blockchain.GetBlockLocator()
		if err != nil {
			log.Errorf("%s Failed to build new header locator: %v", logName, err)
			return
		}

		net.syncCompleted = false
		log.Infof("%s Sending updated locator to peer %s (local best height: %d)", logName, payload.SendFrom, bestHeight)
		net.SendHeaderLocator(payload.SendFrom, NetHeaderLocator{
			SendFrom:   net.Host.ID().String(),
			BestHeight: bestHeight,
			Locator:    locator,
		})
	}

	if len(blocks) == 0 {
		log.Infof("%s No new headers to send to peer %s (best height: %d)", logName, payload.SendFrom, bestHeight)
		net.SendHeaders(payload.SendFrom, NetHeaders{
			SendFrom:   net.Host.ID().String(),
			BestHeight: bestHeight,
			Data:       []NetHeadersData{},
		})
		return
	}

	data := make([]NetHeadersData, 0, len(blocks))
	for _, block := range blocks {
		data = append(data, NetHeadersData{
			Height:   block.Height,
			PrevHash: block.PrevHash,
			Hash:     block.Hash,
			Nbits:    block.NBits,
		})
	}

	log.Infof("%s Sending %d headers to peer %s for synchronization (up to height %d)",
		logName, len(data), payload.SendFrom, bestHeight)

	net.SendHeaders(payload.SendFrom, NetHeaders{
		SendFrom:   net.Host.ID().String(),
		BestHeight: bestHeight,
		Data:       data,
	})

	log.Infof("%s Headers sent successfully", logName)
}

func (net *Network) HandleGetBlockData(content *ChannelContent) {
	const logName = "[GOSSIP::BLOCK]"

	buf := new(bytes.Buffer)
	var payload NetBlock

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("%s Failed to decode block data request from peer %s: %v", logName, content.SendFrom, err)
		return
	}

	block := &blockchain.Block{}
	block = blockchain.DeserializeBlockData(payload.Block)

	log.Infof("%s Received block data %x (height=%d) from peer %s", logName, block.Hash[:6], block.Height, payload.SendFrom)

	hasBlock, err := net.Blockchain.HasBlock(block.Hash)
	if err != nil {
		log.Errorf("%s Error checking block existence %x: %v", logName, block.Hash[:6], err)
		return
	}

	if hasBlock {
		log.Warnf("%s Block %x (height=%d) already exists", logName, block.Hash[:6], block.Height)
		if !net.Gossip.HasSeen(hex.EncodeToString(block.Hash), payload.SendFrom) {
			log.Infof("%s Rebroadcasting header for existing block %x", logName, block.Hash[:6])
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
		log.Errorf("%s Failed to add block %x: %v", logName, block.Hash[:6], err)
		return
	}

	log.Infof("%s Block %d (%x) added to chain", logName, block.Height, block.Hash[:6])

	for _, tx := range block.Transactions {
		txID := hex.EncodeToString(tx.ID)
		MemoryPool.RemoveFromAll(txID)
	}

	if net.Miner && net.IsMining {
		log.Infof("%s Competing block received while mining: %x", logName, block.Hash[:6])
		net.competingBlockChan <- block
	}

	log.Infof("%s Broadcasting block header hash=%x height=%d to peers (excluding sender=%s)", logName, block.Hash[:6], block.Height, payload.SendFrom)
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

	log.Infof("%s Completed handling GetBlockData from peer=%s for block hash=%x height=%d", logName, content.SendFrom, block.Hash[:6], block.Height)
}

func (net *Network) HandleGetData(content *ChannelContent) {
	const logName = "[GOSSIP::DATA]"
	buf := new(bytes.Buffer)
	var payload NetGetData

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("%s Failed to decode GetData request from peer %s: %v", logName, content.SendFrom, err)
		return
	}

	block, err := net.Blockchain.GetBlockMainChain(payload.Hash)
	if err != nil {
		log.Errorf("%s Block %x not found in main chain for peer %s", logName, payload.Hash[:6], payload.SendFrom)
		return
	}

	block = BlockForNetwork(block)
	net.SendFullBlock(payload.SendFrom, &block)

	log.Infof("%s Sent block height=%d hash=%x to peer %s", logName, block.Height, block.Hash[:6], payload.SendFrom)
}

func (net *Network) HandleGetHeader(content *ChannelContent) {
	const logName = "[GOSSIP::HEADER]"

	buf := new(bytes.Buffer)
	var payload NetHeader

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("%s Aborted: failed to decode header request payload", logName)
		return
	}

	exists, err := net.Blockchain.HasBlock(payload.Hash)
	if err != nil {
		log.Errorf("%s Error checking block existence: %v", logName, err)
		return
	}

	lastBlock, err := net.Blockchain.GetLastBlock()
	if err != nil {
		log.Errorf("%s Error retrieving last block: %v", logName, err)
		return
	}

	shortHash := fmt.Sprintf("%x", payload.Hash[:6])
	peerID := payload.SendFrom

	if exists {
		log.Infof("%s Received header for existing block [%s] from peer [%s]", logName, shortHash, peerID)

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
			log.Debugf("%s Broadcasted header [%s] to peers (origin: %s)", logName, shortHash, peerID)
		} else {
			log.Debugf("%s Header [%s] from peer [%s] already seen — ignored", logName, shortHash, peerID)
		}

	} else {
		if bytes.Equal(lastBlock.Hash, payload.PrevHash) {
			net.syncCompleted = true
			log.Infof("%s Header [%s] links to last known block — requesting full block from peer [%s]",
				logName, shortHash, peerID)
			net.SendGetData(peerID, NetGetData{
				SendFrom: net.Host.ID().String(),
				Height:   payload.Height,
				Hash:     payload.Hash,
			})
		} else if payload.Height > lastBlock.Height+1 {
			log.Warnf("%s Peer [%s] appears ahead — header [%s] at height %d vs local %d. Requesting block locator...",
				logName, peerID, shortHash, payload.Height, lastBlock.Height)

			locator, err := net.Blockchain.GetBlockLocator()
			if err != nil {
				log.Errorf("%s Error building block locator: %v", logName, err)
				return
			}

			bestHeight, err := net.Blockchain.GetBestHeight()
			if err != nil {
				log.Errorf("%s Failed to get best height: %v", logName, err)
				return
			}

			net.syncCompleted = false
			net.SendHeaderLocator(peerID, NetHeaderLocator{
				SendFrom:   net.Host.ID().String(),
				BestHeight: bestHeight,
				Locator:    locator,
			})
			log.Infof("%s Sent header locator request to peer [%s]", logName, peerID)
		} else {
			log.Infof("%s Received non-continuous header [%s] from peer [%s] — no action taken",
				logName, shortHash, peerID)
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

func (net *Network) HandleGetRequestGossipPeer(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetRequestGossipPeer

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("[HandleGetRequestGossipPeer] Failed to decode 'NetRequestGossipPeer' request: %v", err)
		return
	}

	if payload.Count > 20 {
		return
	}

	peerList, err := RandomHealthyPeers(int(payload.Count))
	if err != nil {
		log.Errorf("[HandleGetRequestGossipPeer] Failed to get random unique peers: %v", err)
		return
	}

	net.SendGossipPeers(
		payload.SendFrom,
		NetGossipPeers{
			SendFrom: net.Host.ID().String(),
			Peers:    peerList.ListAddrs(),
		},
	)

}

func (net *Network) HandleGetGossipPeers(content *ChannelContent) {
	buf := new(bytes.Buffer)
	var payload NetGossipPeers

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Errorf("[HandleGetGossipPeers] Failed to decode 'NetGossipPeers' request: %v", err)
		return
	}

	ctx := context.Background()

	for _, addr := range payload.Peers {
		log.Infof("Address Peer: %s", addr)
		peerInfo, err := peer.AddrInfoFromString(addr)
		if err != nil {
			log.Errorf("[HandleGetGossipPeers]invalid peer address format: %v", err)
			continue
		}

		if peerInfo.ID == net.Host.ID() {
			continue
		}

		err = AddPeer(addr)
		if err != nil {
			log.Errorf("[HandleGetGossipPeers]Failed to add peer %v", err)
			continue
		}

		err = SafeConnect(ctx, net.Host, addr)
		if err != nil {
			log.Errorf("Failed to connected peer %v", err)
			err := UpdatePeerStatus(addr, false)
			if err != nil {
				log.Errorf("[HandleGetGossipPeers]Failed to update peer %v", err)
			}
			continue
		}

		err = UpdatePeerStatus(addr, true)
		if err != nil {
			log.Errorf("[HandleGetGossipPeers]Failed to update peer %v", err)
		}
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
