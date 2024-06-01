package main

import (
	"context"
	"log"
	"net"
	pp "repMemCache/protoPeer/proto"

	"google.golang.org/grpc"
)

type MasterServer struct {
	pp.UnimplementedPeerServiceServer
	Peers *pp.Peers
}

func (s *MasterServer) RegisterPeer(ctx context.Context, req *pp.RegisterPeerRequest) (*pp.RegisterPeerResponse, error) {
	s.Peers.AddPeer(req.Address)
	log.Printf("Peer is registred %v", req.Address)
	return &pp.RegisterPeerResponse{}, nil
}

func (s *MasterServer) ListPeers(ctx context.Context, req *pp.ListPeersRequest) (*pp.ListPeersResponse, error) {
	peers := s.Peers.GetPeers()
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
	pp.RegisterPeerServiceServer(grpcServer, &MasterServer{Peers: peers})
	log.Printf("Peer Server is running on port: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
