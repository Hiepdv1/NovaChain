package p2p

import (
	"context"
	"core-blockchain/common/helpers"
	"core-blockchain/memopool"
	cryptoRand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
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
	MemoryPool = memopool.Memopool{
		Pending: map[string]memopool.TxInfo{},
		Queued:  map[string]memopool.TxInfo{},
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
	if err != nil {
		log.Error(err)
		return
	}

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
	if err != nil {
		log.Error(err)
		return
	}

	for _, addr := range host.Addrs() {
		log.Infoln("Listening on ", addr)
	}

	log.Info("Host Created: ", host.ID())

	pub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Error(err)
		return
	}

	generalChannel, err := JoinChannel(ctx, pub, host.ID(), GeneralChannel, true)
	if err != nil {
		log.Error(err)
		return
	}

	subscribe := false
	if miner {
		subscribe = true
	}

	miningChannel, err := JoinChannel(ctx, pub, host.ID(), MiningChannel, subscribe)
	if err != nil {
		log.Error(err)
		return
	}

	subscribe = false
	if fullNode || miner {
		subscribe = true
	}

	fullNodesChannel, err := JoinChannel(ctx, pub, host.ID(), FullNodesChannel, subscribe)
	if err != nil {
		log.Error(err)
		return
	}

	ui := NewCLIUI(generalChannel, miningChannel, fullNodesChannel)

	err = SetupDiscovery(ctx, host)
	if err != nil {
		log.Error(err)
		return
	}

	g := NewGossipManager(time.Minute)
	defer g.Stop()
	syncManager := NewSyncManager()

	network := &Network{
		Host:             host,
		MiningChannel:    miningChannel,
		FullNodesChannel: fullNodesChannel,
		Blockchain:       bc,
		Blocks:           make(chan *blockchain.Block, 200),
		Transactions:     make(chan []*blockchain.Transaction, 200),
		Miner:            miner,

		competingBlockChan: make(chan *blockchain.Block, 200),

		peersSyncedWithLocalHeight: []string{},
		syncCompleted:              false,

		Gossip:      g,
		syncManager: syncManager,
	}

	worker := NewWorker(1000, ctx, Error, func(content *ChannelContent) {
		ui.HandleStream(network, content)
	})
	worker.Start(1)

	network.worker = worker

	go HandleEvents(network)

	if miner {
		go network.MinersEventLoop()
	}

	go func() {
		time.Sleep(20 * time.Second)
		network.HandleRequestSync()
	}()

	if err := ui.Run(network, callback); err != nil {
		log.Fatalf("Error running CLI UI: %v", err)
	}

}

func SetupDiscovery(ctx context.Context, host host.Host) error {
	kademliaDHT, err := dht.New(ctx, host, dht.Mode(dht.ModeAuto))
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

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
