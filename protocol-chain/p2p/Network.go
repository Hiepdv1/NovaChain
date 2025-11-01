package p2p

import (
	"context"
	"core-blockchain/common/helpers"
	"core-blockchain/memopool"
	"fmt"
	"math/rand"
	"path"
	"sync"
	"time"

	blockchain "core-blockchain/core"

	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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

	DHT_PREFIX = "/novaChain"
)

var (
	MemoryPool = memopool.Memopool{
		Pending: map[string]memopool.TxInfo{},
		Queued:  map[string]memopool.TxInfo{},
	}
	MinerAddress = ""

	identityFile = path.Join(Root, "/.identity")
)

var bootstrapPeers = []string{
	"/ip4/172.28.0.10/tcp/9000/p2p/Qmb51pbTY5Nu7ERPJQLyK7kMQ96JQTSyPdxyWLcW4Zoq58",
}

func StartNode(bc *blockchain.Blockchain, listenPort, minerAddress string, miner, fullNode, isSeedPeer bool, callback func(*Network)) {

	MinerAddress = minerAddress

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer bc.Database.Close()
	go helpers.CloseDB(bc)

	prvKey, err := LoadOrCreateIdentity(identityFile)
	if err != nil {
		log.Errorf("Failed Load or create identity: %v", err)
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

	ui := NewCLIUI(miningChannel, fullNodesChannel)

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

		Gossip:      g,
		syncManager: syncManager,
	}

	worker := NewWorker(1000, ctx, Error, func(content *ChannelContent) {
		ui.HandleStream(network, content)
	})

	worker.Start(1)
	network.worker = worker

	callback(network)

	go HandleEvents(network)

	err = SetupDiscovery(ctx, host, isSeedPeer)
	if err != nil {
		log.Errorf("Set up discovery failed: %v", err)
		return
	}

	log.Infof("NODE STARTED: Host ID=%s", host.ID())
	log.Infof("NODE ADDRESS LIST (%d):", len(host.Addrs()))
	for _, addr := range host.Addrs() {
		log.Infof(" - %s", addr)
	}

	if isSeedPeer {
		log.Infof("ROLE: Seed Peer ✅ (listening for incoming peers)")
	} else {
		log.Infof("ROLE: Full Node 🌐 (will connect to seed peers for discovery)")
	}

	go network.HandleRequestSync()

	if miner {
		go network.MinersEventLoop()
	}

	if err := ui.Run(network); err != nil {
		log.Errorf("Error running CLI UI: %v", err)
	}

}

func SetupDiscovery(ctx context.Context, host host.Host, isSeedPeer bool) error {
	mode := dht.ModeAuto
	if isSeedPeer {
		mode = dht.ModeServer
	}

	kademliaDHT, err := dht.New(
		ctx,
		host,
		dht.Mode(mode),
		dht.ProtocolPrefix(DHT_PREFIX),
	)
	if err != nil {
		return fmt.Errorf("failed to create DHT: %w", err)
	}

	if isSeedPeer {
		log.Info("🌱 Seed peer initialized. Waiting for incoming peer connections...")
	}

	ConnectBootstrapPeers(ctx, host, kademliaDHT)
	bootstrapDHT(ctx, kademliaDHT)

	go MaintainDHTBootstrap(ctx, kademliaDHT)
	go startDiscoveryTasks(ctx, host, kademliaDHT)
	go RefreshDHT(ctx, host, kademliaDHT)

	return nil
}

func bootstrapDHT(ctx context.Context, dhtInstance *dht.IpfsDHT) {
	if err := dhtInstance.Bootstrap(ctx); err != nil {
		log.Warnf("Initial DHT bootstrap failed: %v", err)
	} else {
		log.Info("✅ DHT bootstrap successful")
	}
}

func startDiscoveryTasks(ctx context.Context, host host.Host, dhtInstance *dht.IpfsDHT) {
	time.Sleep(10 * time.Second)
	go AutoAdvertise(ctx, dhtInstance, Rendezvous)
	go DiscoveryPeers(ctx, host, dhtInstance)
}

func MaintainDHTBootstrap(ctx context.Context, dht *dht.IpfsDHT) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bootstrapDHT(ctx, dht)
		}
	}
}

func ConnectBootstrapPeers(ctx context.Context, host host.Host, kademliaDHT *dht.IpfsDHT) {
	var wg sync.WaitGroup

	for _, addr := range bootstrapPeers {
		peerAddr, err := peer.AddrInfoFromString(addr)
		if err != nil {
			log.Warnf("❌ Invalid bootstrap peer: %v", err)
			continue
		}
		if peerAddr.ID == host.ID() {
			continue
		}

		wg.Add(1)
		go func(p peer.AddrInfo) {
			defer wg.Done()

			if err := host.Connect(ctx, p); err != nil {
				log.Warnf("Failed to connect bootstrap peer %s: %v", p.ID, err)
				return
			}
			log.Infof("✅ Connected bootstrap peer: %s", p.ID)

			tctx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			if _, err := kademliaDHT.FindPeer(tctx, p.ID); err != nil {
				log.Debugf("DHT FindPeer (warmup) returned error for %s: %v", p.ID.String(), err)
			} else {
				log.Infof("🌐 DHT warmed up with peer: %s", p.ID.String())
			}
		}(*peerAddr)
	}

	wg.Wait()
}

func AutoAdvertise(ctx context.Context, dht *dht.IpfsDHT, rendezvous string) {
	routingDiscovery, backoffStrategy := NewOptionsBackoffDiscovery(dht)
	discoveryWithBackoff, err := backoff.NewBackoffDiscovery(routingDiscovery, backoffStrategy)
	if err != nil {
		log.Panicf("❌ Failed to create backoff discovery: %v", err)
		return
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := discoveryWithBackoff.Advertise(ctx, rendezvous); err != nil {
				log.Debugf("Advertise failed (will retry): %v", err)
				continue
			}

			time.Sleep(2 * time.Second)
			peers := dht.RoutingTable().ListPeers()
			if len(peers) == 0 {
				log.Debug("Advertised but routing table still empty")
				continue
			}
			log.Infof("✅ Advertised ourselves successfully to %d peers.", len(peers))
		}
	}
}

func DiscoveryPeers(ctx context.Context, host host.Host, kademliaDHT *dht.IpfsDHT) {
	routingDiscovery, backoffStrategy := NewOptionsBackoffDiscovery(kademliaDHT)
	discoveryWithBackoff, err := backoff.NewBackoffDiscovery(routingDiscovery, backoffStrategy)
	if err != nil {
		log.Panicf("❌ Failed to create backoff discovery: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			peerChan, err := discoveryWithBackoff.FindPeers(ctx, Rendezvous)
			if err != nil {
				log.Warnf("Peer discovery failed: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			for p := range peerChan {
				if p.ID == host.ID() {
					continue
				}
				if host.Network().Connectedness(p.ID) != network.Connected {
					if err := host.Connect(ctx, p); err != nil {
						log.Warnf("Failed to connect discovered peer %s: %v", p.ID, err)
						continue
					}
					log.Infof("✅ Connected discovered peer: %s", p.ID.String())

					kademliaDHT.RoutingTable().PeerAdded(p.ID)
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func RefreshDHT(ctx context.Context, h host.Host, dht *dht.IpfsDHT) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, pid := range h.Network().Peers() {
				_, err := dht.FindPeer(ctx, pid)
				if err == nil {
					log.Infof("[DHT REFRESH] peer %s ok", pid)
				} else {
					log.Warnf("[DHT REFRESH] peer %s not found: %v", pid, err)
				}
			}

			size := len(dht.RoutingTable().ListPeers())
			log.Infof("[DHT REFRESH] routing table size=%d | connected peers=%d",
				size, len(h.Network().Peers()))
		}
	}
}

func NewOptionsBackoffDiscovery(kademliaDHT *dht.IpfsDHT) (*discovery.RoutingDiscovery, backoff.BackoffFactory) {
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	backoffStrategy := backoff.NewExponentialBackoff(
		250*time.Millisecond,
		15*time.Second,
		backoff.FullJitter,
		250*time.Millisecond, 2.0, 50*time.Millisecond,
		rand.NewSource(time.Now().UnixNano()),
	)

	return routingDiscovery, backoffStrategy
}
