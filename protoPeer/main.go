package main

import (
	"context"
	"log"
	"net"
	pp "repMemCache/protoPeer/proto"

	"google.golang.org/grpc"
)

type masterServer struct {
	pp.UnimplementedPeerServiceServer
	peers *pp.Peers
}

func (s *masterServer) RegisterPeer(ctx context.Context, req *pp.RegisterPeerRequest) (*pp.RegisterPeerResponse, error) {
	s.peers.AddPeer(req.Address)
	log.Printf("Peer is registred %v", req.Address)
	return &pp.RegisterPeerResponse{}, nil
}

func (s *masterServer) ListPeers(ctx context.Context, req *pp.ListPeersRequest) (*pp.ListPeersResponse, error) {
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
	peers := pp.NewPeers()
	pp.RegisterPeerServiceServer(grpcServer, &masterServer{peers: peers})
	log.Printf("Peer Server is running on port: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
