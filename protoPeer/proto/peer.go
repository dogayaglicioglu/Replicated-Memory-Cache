package proto

import "sync"

type Peer struct {
	Address string
}

type Peers struct {
	peerList []Peer
	mu       sync.RWMutex
}

func NewPeers() *Peers {
	return &Peers{}
}

func (p *Peers) AddPeer(address string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peerList = append(p.peerList, Peer{Address: address})
}

func (p *Peers) GetPeers() []Peer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.peerList
}
