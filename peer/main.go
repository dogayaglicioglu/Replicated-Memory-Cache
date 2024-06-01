package main

import (
	"context"
	"log"
	"net"
	pp "repMemCache/protoPeer/proto"
	"sync"

	"google.golang.org/grpc"
)

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

type server struct {
	pp.UnimplementedPeerServiceServer
	peers *Peers
}

func (s *server) RegisterPeer(ctx context.Context, req *pp.RegisterPeerRequest) (*pp.RegisterPeerResponse, error) {
	s.peers.AddPeer(req.Address)
	log.Printf("Peer is registred %v", req.Address)
	return &pp.RegisterPeerResponse{}, nil
}

func (s *server) ListPeers(ctx context.Context, req *pp.ListPeersRequest) (*pp.ListPeersResponse, error) {
	peers := s.peers.GetPeers()
	addresses := make([]string, len(peers))
	for i, peer := range peers {
		addresses[i] = peer.Address
	}
	return &pp.ListPeersResponse{Addresses: addresses}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	peers := NewPeers()
	pp.RegisterPeerServiceServer(grpcServer, &server{peers: peers})
	log.Printf("Peer Server is running on port: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
