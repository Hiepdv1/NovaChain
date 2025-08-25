package p2p

import (
	blockchain "core-blockchain/core"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

func BroadcastHeaderRequest(net *Network) {
	peers := net.GeneralChannel.ListPeers()
	if len(peers) == 0 {
		log.Warn("No peers available to broadcast header request")
		net.isSynced <- struct{}{}
		return
	}
	for _, peer := range peers {
		net.SendHeader(peer.String(), true)
	}
}

func (net *Network) SendGetData(sendTo string, data *NetGetData) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_GET_DATA), payload...)
	net.GeneralChannel.Publish("Received getdata", request, sendTo)
}

func (net *Network) SendInv(inv *NetInventory, SendTo string) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_INVENTORY), payload...)
	net.GeneralChannel.Publish("Received Inventory", request, SendTo)
}

func (net *Network) SendGetBlocks(blockHeader *NetGetBlockHeader, SendTo string) {
	payload := GobEncode(blockHeader)
	request := append(CmdToBytes(PREFIX_BLOCKS_HEADER), payload...)

	net.GeneralChannel.Publish("Received getblocks", request, SendTo)
}

func (net *Network) SendTx(sendTo string, transactions []*blockchain.Transaction) {
	serialized := make([][]byte, 0)

	for _, tx := range transactions {
		txID := hex.EncodeToString(tx.ID)
		_, existsPending := memoryPool.Pending[txID]
		_, existsQueued := memoryPool.Queued[txID]

		if !existsPending && !existsQueued {
			memoryPool.Add(*tx)
		}

		serialized = append(serialized, tx.Serializer())
	}

	tnx := NetTx{
		SendFrom:     net.Host.ID().String(),
		Transactions: serialized,
	}

	payload := GobEncode(tnx)
	request := append(CmdToBytes(PREFIX_TX), payload...)
	net.FullNodesChannel.Publish("Received send transaction", request, sendTo)
}

func (net *Network) SendBlock(SendTo string, serializedBlocks [][]byte) {

	data := NetBlock{
		SendFrom: net.Host.ID().String(),
		Blocks:   serializedBlocks,
	}

	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_BLOCK), payload...)

	net.GeneralChannel.Publish("Received send block", request, SendTo)
}

func (net *Network) SendHeader(sendTo string, isSynced bool) {
	lastestBlock, err := net.Blockchain.GetLastBlock()
	if err != nil {
		log.Errorf("Error getting last block: %v", err)
		return
	}
	payload := GobEncode(NetHeader{
		Version:    version,
		BestHeight: lastestBlock.Height,
		Hash:       lastestBlock.Hash,
		SendFrom:   net.Host.ID().String(),
		isSynced:   isSynced,
	})
	request := append(CmdToBytes(PREFIX_HEADER), payload...)
	err = net.GeneralChannel.Publish("Received send header", request, sendTo)
	if err != nil {
		log.Errorf("Error publishing header: %v", err)
		return
	}
}

func (net *Network) SendTxPoolInv(SendTo, _type string, items [][]byte) {
	inv := NetInventory{
		SendFrom: net.Host.ID().String(),
		Type:     _type,
		Items:    items,
	}
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_INVENTORY), payload...)

	net.MiningChannel.Publish("Received Inventory", request, SendTo)
}

func (net *Network) SendTxFromPool(SendTo string, txs []*blockchain.Transaction) {
	var serialize [][]byte

	for _, tx := range txs {
		serialize = append(serialize, tx.Serializer())
	}

	netTx := NetTx{
		SendFrom:     net.Host.ID().String(),
		Transactions: serialize,
	}

	payload := GobEncode(netTx)
	request := append(CmdToBytes(PREFIX_TX), payload...)

	net.MiningChannel.Publish("Received send transactions", request, SendTo)
}
