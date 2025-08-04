package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"core-blockchain/memopool"
	"encoding/gob"
	"encoding/hex"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"
)

func (net *Network) HandleReoganizeTx(txs []*blockchain.Transaction) {
	net.SendTx("", txs)
}

func (net *Network) HandleTx(content *ChannelContent) {
	var buff bytes.Buffer
	var payload NetTx

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding transaction request: ", err)
		return
	}

	if len(payload.Transactions) > 0 {
		chain := net.Blockchain.ContinueBlockchain()
		for _, txBytes := range payload.Transactions {
			newTx := blockchain.Transaction{}
			deserializedTx := newTx.Deserialize(txBytes)
			if chain.VerifyTransaction(deserializedTx) {
				memoryPool.Add(*deserializedTx)
				if net.Miner {
					memoryPool.Move(*deserializedTx, memopool.MEMO_MOVE_FLAG_QUEUED)
					log.Info("MINING")

				}
			}
		}

		net.MineTX(memoryPool.Queued)
	}

	log.Infof("Received Transactions from %s. MemoryPool: %d", content.SendFrom, len(memoryPool.Pending))
}

func (net *Network) HandleGetBlocksData(content *ChannelContent) {
	var buff bytes.Buffer
	var payload NetBlock

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Error("Error decoding block data: ", err)
	}

	for i := len(payload.Blocks) - 1; i > 0; i-- {
		block := blockchain.Block{}

		block = *block.Deserialize(payload.Blocks[i])

		err := net.Blockchain.AddBlock(&block, net.HandleReoganizeTx)

		if err != nil {
			log.Errorf("Error adding block to blockchain: %v", err)
			continue
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			memoryPool.RemoveFromAll(txID)
		}

		if net.Miner {
			log.Info("Competing Block Received From Network: ", block.Hash)
			net.competingBlockChan <- &block
		}

	}

	UTXO := blockchain.UTXOSet{Blockchain: net.Blockchain}
	UTXO.Compute()

	BroadcastHeaderRequest(net)
}

func (net *Network) HandleGetData(content *ChannelContent) {
	var buff bytes.Buffer
	var payload NetGetData

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get data request: ", err)
		return
	}

	if payload.Type == PREFIX_BLOCK {
		serialized := make([][]byte, 0)

		for _, blockHash := range payload.Data {
			block, err := net.Blockchain.GetBlock(blockHash)
			if err != nil {
				log.Errorf("Error getting block from blockchain: %v", err)
				continue
			}

			block = BlockForNetwork(block)

			serialized = append(serialized, block.Serialize())
		}

		net.SendBlock(payload.SendFrom, serialized)

	}

	if payload.Type == PREFIX_TX {
		txs := []*blockchain.Transaction{}

		if net.BelongsToMiningGroup(payload.SendFrom) {
			for _, txBytes := range payload.Data {
				txID := hex.EncodeToString(txBytes)
				tx := memoryPool.Pending[txID]
				memoryPool.Move(tx, memopool.MEMO_MOVE_FLAG_QUEUED)
			}
			net.SendTxFromPool(payload.SendFrom, txs)
		} else {
			for _, txBytes := range payload.Data {
				txID := hex.EncodeToString(txBytes)
				tx := memoryPool.Pending[txID]
				txs = append(txs, &tx)
			}
			net.SendTx(payload.SendFrom, txs)
		}
	}
}

func (net *Network) HandleGetInventory(content *ChannelContent) {
	var buff bytes.Buffer
	var inv NetInventory

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&inv)
	if err != nil {
		log.Error("Error decoding inventory request: ", err)
		return
	}

	log.Infof("Received inventory with %d %s \n", len(inv.Items), inv.Type)

	if inv.Type == PREFIX_BLOCK {
		if len(inv.Items) > 0 {
			net.SendGetData(inv.SendFrom, &NetGetData{
				SendFrom: net.Host.ID().String(),
				Type:     inv.Type,
				Data:     inv.Items,
			})
		} else {
			log.Info("Empty block hashes")
		}
	}

	if inv.Type == PREFIX_TX {
		txIDs := [][]byte{}
		for _, txID := range inv.Items {
			if memoryPool.Pending[hex.EncodeToString(txID)].ID == nil {
				txIDs = append(txIDs, txID)
			}
		}
		if len(txIDs) > 0 {
			netData := NetGetData{
				SendFrom: net.Host.ID().String(),
				Type:     PREFIX_TX,
				Data:     txIDs,
			}
			net.SendGetData(inv.SendFrom, &netData)
		}
	}
}

func (net *Network) HandleGetBlocksHeader(content *ChannelContent) {
	var buff bytes.Buffer
	var payload *NetGetBlockHeader

	buff.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get blocks header request: ", err)
		return
	}

	chain := net.Blockchain.ContinueBlockchain()
	blockHashes := chain.GetBlockHashes(payload.Hash, payload.Height, MAX_BLOCKS_HEADER)
	log.Info("LENGTH: ", len(blockHashes), " MAX: ", MAX_BLOCKS_HEADER)
	net.SendInv(&NetInventory{
		SendFrom: net.Host.ID().String(),
		Type:     PREFIX_BLOCK,
		Items:    blockHashes,
	}, payload.SendFrom)

}

func (net *Network) HandleGetHeader(content *ChannelContent) {
	var buf bytes.Buffer
	var payload NetHeader

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding header request: ", err)
		return
	}

	bestHeight := net.Blockchain.GetBestHeight()
	otherHeight := payload.BestHeight

	log.Info("BEST HEIGHT: ", bestHeight, " OTHER HEIGHT: ", otherHeight)
	if bestHeight < otherHeight {
		net.SendGetBlocks(&NetGetBlockHeader{
			SendFrom: net.Host.ID().String(),
			Height:   bestHeight,
			Hash:     payload.Hash,
		}, payload.SendFrom)
	} else if bestHeight > otherHeight {
		net.SendHeader(payload.SendFrom)
	} else {
		if !net.syncCompleted {
			if len(net.peersSyncedWithLocalHeight) == len(net.GeneralChannel.ListPeers()) {
				log.Info("All peers synced with local height, marking network as synced")
				net.isSynced <- struct{}{}
			} else if !slices.Contains(net.peersSyncedWithLocalHeight, payload.SendFrom) {
				log.Info("Waiting for more peers to sync with local height")
				net.peersSyncedWithLocalHeight = append(net.peersSyncedWithLocalHeight, payload.SendFrom)
			}
		}
	}
}

func (net *Network) HandleGetTxFromPool(content *ChannelContent) {
	var buf bytes.Buffer
	var payload *TxFromPool

	buf.Write(content.Payload[commandLength:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Error("Error decoding get transaction from pool request: ", err)
		return
	}

	if int64(len(memoryPool.Pending)) >= payload.Count {
		txs := memoryPool.GetTransactions(payload.Count)
		net.SendTxPoolInv(payload.SendFrom, PREFIX_TX, txs)
	} else {
		net.SendTxPoolInv(payload.SendFrom, PREFIX_TX, [][]byte{})
	}
}

func HandleEvents(net *Network) {
	for {
		select {
		case block := <-net.Blocks:
			net.SendBlock("", [][]byte{block.Serialize()})
		case tx := <-net.Transactions:
			net.SendTx("", []*blockchain.Transaction{tx})
		}
	}
}

func (net *Network) MinersEventLoop() {
	poolCheckTicker := time.NewTicker(time.Second)
	defer poolCheckTicker.Stop()

	for range poolCheckTicker.C {
		tnx := TxFromPool{
			SendFrom: net.Host.ID().String(),
			Count:    1,
		}
		payload := GobEncode(tnx)
		request := append(CmdToBytes(PREFIX_TX_FROM_POOL), payload...)
		net.FullNodesChannel.Publish("Request transaction from pool", request, "")
	}
}

func (net *Network) MineTX(memopoolTxs map[string]blockchain.Transaction) {
	var txs []*blockchain.Transaction
	log.Infof("MINE: %d", len(memopoolTxs))
	chain := net.Blockchain.ContinueBlockchain()

	for id := range memopoolTxs {
		tx := memopoolTxs[id]
		log.Infof("TX: %s \n", tx.ID)

		if chain.VerifyTransaction(&tx) {
			log.Info("Valid Transaction")
			txs = append(txs, &tx)
		} else {
			log.Info("Invalid Transaction")
		}
	}

	if len(txs) == 0 {
		log.Error("No Valid Transaction")
		return
	}

	lastCurrentBlock, err := chain.GetLastBlock()
	if err != nil {
		log.Errorf("Error getting last block: %v", err)
		return
	}
	log.Info("Last Block Height: ", lastCurrentBlock.Height)

	newBlock := make(chan *blockchain.Block, 1)
	stopMining := make(chan bool, 1)

	go func() {
		for block := range net.competingBlockChan {
			log.Info("Competing Block Received: ", block.Hash)
			currentChainHead, err := chain.GetLastBlock()
			if err != nil {
				log.Errorf("Error getting last block: %v", err)
				continue
			}
			if block.Height >= lastCurrentBlock.Height+1 && bytes.Equal(block.Hash, currentChainHead.Hash) {
				select {
				case stopMining <- true:
					log.Info("Competing Block is newer or same height, signalling to stop mining.")
				default:
					log.Debug("Stop signal already sent or channel full.")
				}
				return
			} else {
				log.Info("Competing Block is older or irrelevant, ignoring for now.")
			}
		}
	}()

	go func() {
		block, err := chain.MineBlock(txs, MinerAddress, net.HandleReoganizeTx)
		if err != nil {
			log.Errorf("Error mining block: %v", err)
			newBlock <- nil
			return
		}
		newBlock <- block
	}()

	select {
	case <-stopMining:
		log.Info("Mining stopped due to competing block")
		return
	case minedBlock := <-newBlock:
		if minedBlock == nil {
			log.Info("No block mined, stopping mining process.")
			return
		}

		UTXOs := blockchain.UTXOSet{
			Blockchain: chain,
		}
		UTXOs.Compute()

		log.Info("New Block Mined")

		var hashes [][]byte
		hashes = append(hashes, minedBlock.Hash)

		inv := NetInventory{
			SendFrom: net.Host.ID().String(),
			Type:     PREFIX_BLOCK,
			Items:    hashes,
		}

		net.SendInv(&inv, "")

		memoryPool.ClearAll()
	}
}
