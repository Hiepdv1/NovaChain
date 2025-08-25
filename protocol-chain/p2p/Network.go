package p2p

import (
	"context"
	"core-blockchain/common/helpers"
	"core-blockchain/common/utils"
	"core-blockchain/memopool"
	cryptoRand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"time"

	blockchain "core-blockchain/core"

	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/backoff"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	yamux "github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"

	tcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	GeneralChannel   = "general-channel"
	MiningChannel    = "mining-channel"
	FullNodesChannel = "fullnodes-channel"
	Rendezvous       = "room-chain"

	version       = 1
	commandLength = 20
)

var (
	memoryPool = memopool.Memopool{
		Pending: map[string]blockchain.Transaction{},
		Queued:  map[string]blockchain.Transaction{},
	}
	MinerAddress = ""
)

func StartNode(bc *blockchain.Blockchain, listenPort, minerAddress string, miner, fullNode bool, callback func(*Network)) {
	var r io.Reader = cryptoRand.Reader

	MinerAddress = minerAddress

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer bc.Database.Close()
	go helpers.CloseDB(bc)

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	utils.ErrorHandle(err)

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer(yamux.ID, yamux.DefaultTransport),
		libp2p.Muxer(mplex.ID, mplex.DefaultTransport),
	)

	listenAddr := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", listenPort),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", listenPort),
	)

	host, err := libp2p.New(
		transports,
		muxers,
		listenAddr,

		libp2p.Identity(prvKey),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.EnableRelay(),
		libp2p.EnableRelayService(),
		libp2p.Security(noise.ID, noise.New),
	)
	utils.ErrorHandle(err)

	for _, addr := range host.Addrs() {
		log.Infoln("Listening on ", addr)
	}

	log.Info("Host Created: ", host.ID())

	pub, err := pubsub.NewGossipSub(ctx, host)
	utils.ErrorHandle(err)

	generalChannel, err := JoinChannel(ctx, pub, host.ID(), GeneralChannel, true)
	utils.ErrorHandle(err)

	subscribe := false
	if miner {
		subscribe = true
	}

	miningChannel, err := JoinChannel(ctx, pub, host.ID(), MiningChannel, subscribe)
	utils.ErrorHandle(err)

	subscribe = false
	if fullNode {
		subscribe = true
	}

	fullNodesChannel, err := JoinChannel(ctx, pub, host.ID(), FullNodesChannel, subscribe)
	utils.ErrorHandle(err)

	ui := NewCLIUI(generalChannel, miningChannel, fullNodesChannel)

	err = SetupDiscovery(ctx, host)
	utils.ErrorHandle(err)

	network := &Network{
		Host:             host,
		GeneralChannel:   generalChannel,
		MiningChannel:    miningChannel,
		FullNodesChannel: fullNodesChannel,
		Blockchain:       bc,
		Blocks:           make(chan *blockchain.Block, 200),
		Transactions:     make(chan *blockchain.Transaction, 200),
		Miner:            miner,

		competingBlockChan:         make(chan *blockchain.Block, 200),
		blockProcessingLimiter:     make(chan struct{}, 200),
		txProcessingLimiter:        make(chan struct{}, runtime.NumCPU()),
		peersSyncedWithLocalHeight: []string{},
		isSynced:                   make(chan struct{}, 1),
		syncCompleted:              false,
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)

	EventLoop:
		for {
			select {
			case <-ticker.C:
				log.Info("Broadcasting header request to all peers")
				BroadcastHeaderRequest(network)
			case <-network.isSynced:
				log.Info("Network is synced with local height")
				network.syncCompleted = true
				break EventLoop
			}
		}

		go HandleEvents(network)

		if miner {
			go network.MinersEventLoop()
		}
	}()

	if err := ui.Run(network, callback); err != nil {
		log.Fatalf("Error running CLI UI: %v", err)
	}

}

func SetupDiscovery(ctx context.Context, host host.Host) error {
	kademliaDHT, err := dht.New(ctx, host, dht.Mode(dht.ModeAuto))
	utils.ErrorHandle(err)

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerInfo); err != nil {
				log.Errorln(err)
			} else {
				log.Info("Connection established with bootstrap node: ", *peerInfo)
			}
		}()
	}
	wg.Wait()

	log.Info("Bootstrapping the DHT")
	if err := kademliaDHT.Bootstrap(ctx); err != nil {
		log.Panic(err)
	}

	log.Info("Annoucing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	backoffStrategy := backoff.NewExponentialBackoff(time.Second,
		3*time.Minute,
		backoff.FullJitter,
		time.Second,
		2.0,
		100*time.Millisecond,
		rand.NewSource(time.Now().UnixNano()),
	)

	discoveryWithBackoff, err := backoff.NewBackoffDiscovery(
		routingDiscovery,
		backoffStrategy,
	)
	if err != nil {
		log.Fatalf("Failed to create backoff discovery: %v", err)
	}

	_, err = discoveryWithBackoff.Advertise(ctx, Rendezvous)
	if err != nil {
		log.Fatal("Failed to announce ourselves: ", err)
	}
	log.Info("Successfully annouced!")

	log.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, Rendezvous)
	utils.ErrorHandle(err)

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}

		log.Info("Found peer: ", peer)

		log.Debugf("🔗 Connecting to peer %s at %v\n", peer.ID, peer.Addrs)
		err := host.Connect(ctx, peer)
		if err != nil {
			log.Warningf("⚠️ Error connecting to peer %s: %s\n", peer.ID, err)
		}
	}

	return nil
}

func (net *Network) BelongsToMiningGroup(peerID string) bool {
	peers := net.MiningChannel.ListPeers()

	for _, peer := range peers {
		ID := peer.String()

		if ID == peerID {
			return true
		}
	}

	return false
}
