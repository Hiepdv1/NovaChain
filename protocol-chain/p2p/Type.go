package p2p

import (
	blockchain "core-blockchain/core"

	"github.com/libp2p/go-libp2p/core/host"
)

type Network struct {
	Host             host.Host
	GeneralChannel   *Channel
	MiningChannel    *Channel
	FullNodesChannel *Channel
	Blockchain       *blockchain.Blockchain
	Blocks           chan *blockchain.Block
	Transactions     chan *blockchain.Transaction
	Miner            bool

	competingBlockChan         chan *blockchain.Block
	blockProcessingLimiter     chan struct{}
	txProcessingLimiter        chan struct{}
	peersSyncedWithLocalHeight []string
	syncCompleted              bool
	isSynced                   chan struct{}
}

type NetHeader struct {
	Version    int64
	BestHeight int64
	Hash       []byte
	SendFrom   string
	isSynced   bool
}

type NetGetBlockHeader struct {
	SendFrom string
	Height   int64
	Hash     []byte
}

type NetBlock struct {
	SendFrom string
	Blocks   [][]byte
}

type NetTx struct {
	SendFrom     string
	Transactions [][]byte
}

type TxFromPool struct {
	SendFrom string
	Count    int64
}

type NetInventory struct {
	SendFrom string
	Type     string
	Items    [][]byte
}

type NetGetData struct {
	SendFrom string
	Type     string
	Data     [][]byte
}
