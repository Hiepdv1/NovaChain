package p2p

import (
	"bytes"
	blockchain "core-blockchain/core"
	"core-blockchain/memopool"

	log "github.com/sirupsen/logrus"
)

func (net *Network) SendTx(sendTo string, tx *blockchain.Transaction) {
	txInfo := memopool.GetTxInfo(tx, net.Blockchain)
	if txInfo == nil {
		return
	}

	MemoryPool.Add(*txInfo)

	buf := new(bytes.Buffer)

	blockchain.SerializeTransaction(tx, buf)

	tnx := NetTx{
		SendFrom:    net.Host.ID().String(),
		Transaction: buf.Bytes(),
	}

	payload := GobEncode(tnx)
	request := append(CmdToBytes(PREFIX_TX), payload...)
	net.FullNodesChannel.Publish("Sending transaction", request, sendTo)
}

func (net *Network) SendFullBlock(sendTo string, block *blockchain.Block) {

	serialize := blockchain.SerializeBlock(block)

	data := NetBlock{
		SendFrom: net.Host.ID().String(),
		Block:    serialize,
	}

	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_BLOCK), payload...)
	net.FullNodesChannel.Publish("Sending full block", request, sendTo)
}

func (net *Network) SendBlockDataSync(sendTo string, data NetBlockSync) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_BLOCK_SYNC), payload...)
	err := net.FullNodesChannel.Publish("Sending block data (sync)", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish block data sync: %v", err)
		return
	}
}

func (net *Network) SendHeaderLocator(sendTo string, data NetHeaderLocator) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_HEADER_LOCATOR), payload...)
	err := net.FullNodesChannel.Publish("Sending header locator", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish header locator: %v", err)
		return
	}
}

func (net *Network) SendHeaders(sendTo string, data NetHeaders) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_HEADER_SYNC), payload...)
	err := net.FullNodesChannel.Publish("Sending headers", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish headers: %v", err)
		return
	}
}

func (net *Network) SendHeader(sendTo string, data *NetHeader) {
	payload := GobEncode(data)

	request := append(CmdToBytes(PREFIX_HEADER), payload...)
	err := net.FullNodesChannel.Publish("Sending single header", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish header: %v", err)
		return
	}
}

func (net *Network) SendGetDataSync(sendTo string, data NetGetDataSync) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_GET_DATA_SYNC), payload...)
	err := net.FullNodesChannel.Publish("Requesting block data (sync)", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish GetDataSync request: %v", err)
		return
	}
}

func (net *Network) SendGetData(sendTo string, data NetGetData) {
	payload := GobEncode(data)
	request := append(CmdToBytes(PREFIX_GET_DATA), payload...)
	err := net.FullNodesChannel.Publish("Requesting block data", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish GetData request: %v", err)
		return
	}
}

func (net *Network) SendTransactions(sendTo string, inv NetTransactionData) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_TXS_Data), payload...)
	err := net.FullNodesChannel.Publish("Sending multiple transactions", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish transactions: %v", err)
		return
	}
}

func (net *Network) SendGetDataTransaction(sendTo string, inv NetGetDataTransaction) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_DATA_TX), payload...)
	err := net.FullNodesChannel.Publish("Requesting specific transactions", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish GetDataTransaction request: %v", err)
		return
	}
}

func (net *Network) SendRequestTxFromPool(sendTo string, inv TxFromPool) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_TX_FROM_POOL), payload...)
	err := net.FullNodesChannel.Publish("Requesting transaction from mempool", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish TxFromPool request: %v", err)
		return
	}
}

func (net *Network) SendTxPoolInv(sendTo string, txHashes [][]byte) {
	inv := NetInventoryTxs{
		SendFrom: net.Host.ID().String(),
		TxHashes: txHashes,
	}
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_INVENTORY), payload...)
	err := net.MiningChannel.Publish("Sending mempool inventory", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish TxPool inventory: %v", err)
		return
	}
}

func (net *Network) SendRequestGossipPeer(sendTo string, inv NetRequestGossipPeer) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFIX_REQUEST_GOSSIP_PEERS), payload...)
	err := net.FullNodesChannel.Publish("Sending Request gossip peers", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish Request gossip peers: %v", err)
		return
	}
}

func (net *Network) SendGossipPeers(sendTo string, inv NetGossipPeers) {
	payload := GobEncode(inv)
	request := append(CmdToBytes(PREFOX_GOSSIP_PEERS), payload...)
	err := net.FullNodesChannel.Publish("Sending gossip peers", request, sendTo)
	if err != nil {
		log.Errorf("Failed to publish gossip peers: %v", err)
		return
	}
}
