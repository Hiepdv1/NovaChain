package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	log "github.com/sirupsen/logrus"
)

var (
	peersFile = path.Join(Root, "/.chain/known_peers.json")
	peerLock  sync.Mutex

	PeerMaxAge     = 7 * 24 * time.Hour
	PeerCheckCycle = 1 * time.Hour
	MaxFailCount   = 3
)

type PeerEntry struct {
	Address      string
	LastSeen     time.Time
	SuccessCount int
	FailCount    int
}

type PeerList struct {
	Peers []PeerEntry
}

func (p *PeerList) ListAddrs() []string {
	addrs := make([]string, 0)
	for _, peerEntry := range p.Peers {
		addrs = append(addrs, peerEntry.Address)
	}

	return addrs
}

func SyncConnectedPeers(ctx context.Context, h host.Host) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers := h.Network().Peers()
			for _, pid := range peers {
				addrs := GetPeerMultiaddr(h, pid)
				for _, addr := range addrs {
					_ = AddPeer(addr)
				}
			}
		}
	}
}

func GetPeerMultiaddr(h host.Host, pid peer.ID) []string {
	peerInfo := h.Peerstore().PeerInfo(pid)
	addrs := []string{}

	for _, addr := range peerInfo.Addrs {
		protocols := addr.Protocols()
		if len(protocols) < 2 {
			continue
		}

		if protocols[0].Name == "ip4" && protocols[1].Name == "tcp" {
			ip, err := addr.ValueForProtocol(4)
			if err != nil {
				log.Printf("Can't get IP from multiaddr: %v", err)
				continue
			}

			if ip == "127.0.0.1" || ip == "0.0.0.0" {
				continue
			}

			fullAddr := fmt.Sprintf("%s/p2p/%s", addr.String(), pid.String())
			addrs = append(addrs, fullAddr)
		}
	}

	return addrs
}

func ensurePeerFile() error {
	if _, err := os.Stat(peersFile); os.IsNotExist(err) {
		list := PeerList{Peers: []PeerEntry{}}
		data, _ := json.MarshalIndent(list, "", "")
		if err := os.MkdirAll(path.Dir(peersFile), 0o755); err != nil {
			return err
		}
		return os.WriteFile(peersFile, data, 0o644)
	}
	return nil
}

func LoadPeers() (*PeerList, error) {
	if err := ensurePeerFile(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(peersFile)
	if err != nil {
		return nil, err
	}

	var list PeerList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

func SavePeers(list *PeerList) error {
	peerLock.Lock()
	defer peerLock.Unlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(peersFile, data, 0o644)
}

func AddPeer(addr string) error {
	list, err := LoadPeers()
	if err != nil {
		return err
	}

	now := time.Now()
	for i, p := range list.Peers {
		if p.Address == addr {
			list.Peers[i].LastSeen = now
			return SavePeers(list)
		}
	}

	entry := PeerEntry{
		Address:      addr,
		LastSeen:     now,
		SuccessCount: 0,
		FailCount:    0,
	}

	list.Peers = append(list.Peers, entry)
	return SavePeers(list)
}

func UpdatePeerStatus(addr string, success bool) error {
	list, err := LoadPeers()
	if err != nil {
		return err
	}

	isExist := false

	for i, p := range list.Peers {
		if p.Address == addr {
			if success {
				list.Peers[i].SuccessCount++
				list.Peers[i].LastSeen = time.Now()
			} else {
				list.Peers[i].FailCount++
			}
			isExist = true
			break
		}
	}

	if !isExist {
		AddPeer(addr)
		UpdatePeerStatus(addr, success)
	}

	return SavePeers(list)
}

func LoadPeerInfos() ([]peer.AddrInfo, error) {
	list, err := LoadPeers()
	if err != nil {
		return nil, err
	}

	infos := make([]peer.AddrInfo, 0, len(list.Peers))
	for _, peerEntry := range list.Peers {
		pi, err := peer.AddrInfoFromString(peerEntry.Address)
		if err != nil {
			fmt.Printf("❌ Invalid peer addr in file: %s\n", peerEntry.Address)
			continue
		}
		infos = append(infos, *pi)
	}
	return infos, nil
}

func SafeConnect(ctx context.Context, h host.Host, addr string) error {
	peerInfo, err := peer.AddrInfoFromString(addr)
	if err != nil {
		return fmt.Errorf("invalid peer address format: %w", err)
	}

	pid := peerInfo.ID

	if pid == h.ID() {
		return fmt.Errorf("skip self-connection: cannot connect to itself")
	}

	if h.Network().Connectedness(pid) == network.Connected {
		log.Infof("[INFO] Already connected to peer %s", pid)
		return nil
	}

	if err := h.Connect(ctx, *peerInfo); err != nil {
		return fmt.Errorf("failed to connect to peer %s: %w", pid, err)
	}

	log.Printf("[OK] Successfully connected to peer %s", pid)
	return nil
}

func RandomHealthyPeers(n int) (*PeerList, error) {
	list, err := LoadPeers()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	valid := []PeerEntry{}
	for _, p := range list.Peers {
		if p.FailCount < MaxFailCount && now.Sub(p.LastSeen) < PeerMaxAge {
			valid = append(valid, p)
		}
	}

	if n >= len(valid) {
		return &PeerList{
			Peers: valid,
		}, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(valid), func(i, j int) {
		valid[i], valid[j] = valid[j], valid[i]
	})

	return &PeerList{
		Peers: valid[:n],
	}, nil
}

func PingPeer(ctx context.Context, h host.Host, addr string) bool {
	pi, err := peer.AddrInfoFromString(addr)
	if err != nil {
		return false
	}

	if pi.ID == h.ID() {
		return false
	}

	pinger := ping.NewPingService(h)

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := pinger.Ping(ctxPing, pi.ID)

	select {
	case res := <-ch:
		if res.Error != nil {
			log.Warnf("Ping failed for %s: %v", pi.ID, res.Error)
			return false
		}
		log.Infof("Ping successful for %s: RTT=%s", pi.ID, res.RTT)
		return true
	case <-ctxPing.Done():
		log.Warnf("Ping timeout for %s", pi.ID)
		return false
	}
}

func CleanupAndPingPeers(ctx context.Context, h host.Host, maxAge time.Duration) {
	list, err := LoadPeers()
	if err != nil {
		log.Warnf("❌ Failed to load peers for cleanup: %v", err)
		return
	}

	cutoff := time.Now().Add(-maxAge)
	active := []PeerEntry{}
	for _, p := range list.Peers {
		if p.LastSeen.Before(cutoff) {
			log.Infof("🧹 Removing expired peer: %s", p.Address)
			continue
		}

		if p.FailCount >= 3 {
			log.Infof("🔍 Checking stale peer: %s", p.Address)
			if PingPeer(ctx, h, p.Address) {
				p.FailCount = 0
				p.LastSeen = time.Now()
				log.Infof("✅ Peer revived: %s", p.Address)
			} else {
				log.Warnf("❌ Peer still unreachable: %s", p.Address)
			}
		}

		active = append(active, p)
	}

	list.Peers = active
	err = SavePeers(list)
	if err != nil {
		log.Errorf("[CleanupAndPingPeers]Failed save peers %v", err)
	}
}

func StartPeerMaintenance(ctx context.Context, h host.Host) {
	ticker := time.NewTicker(PeerCheckCycle)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			CleanupAndPingPeers(ctx, h, PeerMaxAge)
		case <-ctx.Done():
			return
		}
	}
}
