package p2p

import (
	"core-blockchain/common/env"
	"math/big"
	"sync"
	"time"
)

type PeerStatus struct {
	ID        string
	Height    int64
	TotalWork *big.Int
	IsTarget  bool
	LastSeen  time.Time
}

var conf = env.New()

type SyncManager struct {
	mu     sync.Mutex
	peers  map[string]*PeerStatus
	target *PeerStatus
}

func NewSyncManager() *SyncManager {
	syncManager := &SyncManager{
		peers: make(map[string]*PeerStatus),
	}

	go syncManager.RemoveStalePeersLoop(time.Duration(conf.Peer_TTL_Minute) * time.Minute)

	return syncManager
}

func (sm *SyncManager) UpdatePeerStatus(pID string, height int64, work *big.Int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	p, ok := sm.peers[pID]
	if !ok {
		p = &PeerStatus{ID: pID}
		sm.peers[pID] = p
	}
	p.Height = height
	p.TotalWork = new(big.Int).Set(work)
	p.LastSeen = time.Now()

	sm.selectBestPeerLocked()
}

func (sm *SyncManager) selectBestPeerLocked() {
	var best *PeerStatus = sm.target
	for _, p := range sm.peers {
		if best == nil {
			best = p
			continue
		}
		cmp := p.TotalWork.Cmp(best.TotalWork)
		if cmp > 0 || (cmp == 0 && p.Height > best.Height) {
			best = p
		}
	}

	if best != nil && (sm.target == nil || sm.target.ID != best.ID) {
		if sm.target != nil {
			sm.target.IsTarget = false
		}
		best.IsTarget = true
		sm.target = best
	}
}

func (sm *SyncManager) GetTargetPeer() *PeerStatus {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.target
}

func (sm *SyncManager) IsSynced(localHeight int64, localWork *big.Int) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.target == nil {
		return false
	}

	return localHeight >= sm.target.Height && localWork.Cmp(sm.target.TotalWork) >= 0
}

func (sm *SyncManager) ClearTarget() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.target != nil {
		sm.target.IsTarget = false
		sm.target = nil
	}
}

func (sm *SyncManager) RemoveStalePeersLoop(timeout time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now()

	for _, p := range sm.peers {
		if now.Sub(p.LastSeen) > timeout {
			if sm.target != nil && sm.target.ID == p.ID {
				sm.target = nil
			}
			delete(sm.peers, p.ID)
		}
	}

	if sm.target == nil && len(sm.peers) > 0 {
		sm.selectBestPeerLocked()
	}
}
