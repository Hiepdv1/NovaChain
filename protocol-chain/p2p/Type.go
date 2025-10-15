package p2p

import (
	blockchain "core-blockchain/core"

	"github.com/libp2p/go-libp2p/core/host"
)

type Network struct {
	Host             host.Host
	MiningChannel    *Channel
	FullNodesChannel *Channel
	Blockchain       *blockchain.Blockchain
	Blocks           chan *blockchain.Block
	Transactions     chan []*blockchain.Transaction
	Miner            bool

	IsMining                   bool
	competingBlockChan         chan *blockchain.Block
	peersSyncedWithLocalHeight []string
	syncCompleted              bool

	// cache block/transaction - Gossip
	Gossip      *GossipManager
	syncManager *SyncManager

	worker *Worker[*ChannelContent]
}

type NetHeadersData struct {
	Height   int64
	PrevHash []byte
	Hash     []byte
	Nbits    uint32
}

type NetHeaders struct {
	SendFrom   string
	BestHeight int64
	Data       []NetHeadersData
}

type NetBlockSync struct {
	SendFrom   string
	BestHeight int64
	Blocks     []blockchain.Block
}

type NetHeader struct {
	Hash     []byte
	Height   int64
	PrevHash []byte
	SendFrom string
}

type NetHeaderLocator struct {
	SendFrom string
	Locator  [][]byte
}

type NetGetDataSync struct {
	SendFrom string
	Hashes   [][]byte
}

type NetGetData struct {
	SendFrom string
	Height   int64
	Hash     []byte
}

type NetTxMining struct {
	Txs []blockchain.Transaction
}

type NetBlock struct {
	SendFrom string
	Block    []byte
}

type TxFromPool struct {
	SendFrom string
	Count    int64
}

type NetTx struct {
	SendFrom    string
	Transaction []byte
}

type NetInventoryTxs struct {
	SendFrom string
	TxHashes [][]byte
	Count    uint64
}

type NetGetDataTransaction struct {
	SendFrom string
	TxHashes [][]byte
}

type NetTransactionData struct {
	SendFrom     string
	Transactions []blockchain.Transaction
}
