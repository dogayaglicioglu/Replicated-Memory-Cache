package main

import (
	"context"
	"log"
	"net"
	cacheProto "repMemCache/cache/proto"
	pp "repMemCache/protoPeer/proto"
	"sync"

	"google.golang.org/grpc"
)

type MasterServer struct {
	pp.UnimplementedPeerServiceServer
	cacheProto.UnimplementedCacheServiceServer
	Peers *pp.Peers
	mu    sync.Mutex
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

func (s *MasterServer) SetData(ctx context.Context, req *cacheProto.DataRequest) (*cacheProto.DataResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, peer := range s.Peers.GetPeers() {
		conn, err := grpc.Dial(peer.Address, grpc.WithInsecure())
		if err != nil {
			log.Printf("Error in connecting to the peer %v", err)
			continue
		}
		defer conn.Close()
		cacheServiceClient := cacheProto.NewCacheServiceClient(conn)
		_, err = cacheServiceClient.SetData(context.Background(), req)
		if err != nil {
			log.Printf("Error in replicating data to peer %v", err)
		}
		client := pp.NewPeerServiceClient(conn)
		_, err = client.NotifyPeers(context.Background(), &pp.NotifyPeerRequest{
			Message: "Replicated data: " + req.Key + " = " + req.Value,
		})

		if err != nil {
			log.Printf("Error notifying peer: %v", err)
		}
	}

	return &cacheProto.DataResponse{Value: req.Value}, nil
}

func (s *MasterServer) GetData(ctx context.Context, req *cacheProto.KeyRequest) (*cacheProto.DataResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("GetData in master: %v", req.Key)
	return &cacheProto.DataResponse{Value: "example_value"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	peers := pp.NewPeers()
	ms := &MasterServer{Peers: peers}
	pp.RegisterPeerServiceServer(grpcServer, ms)
	cacheProto.RegisterCacheServiceServer(grpcServer, ms)
	log.Printf("Peer Server is running on port: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
