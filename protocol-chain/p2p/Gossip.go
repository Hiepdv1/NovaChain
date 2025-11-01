package p2p

import (
	"slices"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type GossipManager struct {
	mu   sync.Mutex
	seen map[string]map[string]time.Time // hash -> peerID -> timeSeen
	ttl  time.Duration
	quit chan struct{}
}

func NewGossipManager(ttl time.Duration) *GossipManager {
	g := &GossipManager{
		seen: make(map[string]map[string]time.Time),
		ttl:  ttl,
		quit: make(chan struct{}),
	}
	go g.cleanupLoop()
	return g
}

func (g *GossipManager) MarkSeen(hash, peerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.seen[hash]; !ok {
		g.seen[hash] = make(map[string]time.Time)
	}

	g.seen[hash][peerID] = time.Now()
}

func (g *GossipManager) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.cleanup()
		case <-g.quit:
			return
		}

	}
}

func (g *GossipManager) cleanup() {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()
	for hash, peers := range g.seen {
		for pId, t := range peers {
			if now.Sub(t) >= g.ttl {
				delete(peers, pId)
			}
		}
		if len(peers) == 0 {
			delete(g.seen, hash)
		}
	}
}

func (g *GossipManager) HasSeen(hash, peerID string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if peers, ok := g.seen[hash]; ok {
		if _, ok := peers[peerID]; ok {
			return true
		}
	}
	return false
}

func (g *GossipManager) Stop() {
	close(g.quit)
}

func (g *GossipManager) Broadcast(
	peers []peer.ID,
	excludePeerIDs []string,
	handler func(p peer.ID),
) {
	if len(peers) < 1 {
		return
	}

	for _, pId := range peers {
		if slices.Contains(excludePeerIDs, pId.String()) {
			continue
		}
		handler(pId)
	}
}
