package cache

type Peer struct {
	Address string
}

type Peers struct {
	peerList []Peer
}

func NewPeers() *Peers {
	return &Peers{}
}

func (p *Peers) AddPeer(address string) {
	p.peerList = append(p.peerList, Peer{Address: address})
}

func (p *Peers) GetPeers() []Peer {
	return p.peerList
}
