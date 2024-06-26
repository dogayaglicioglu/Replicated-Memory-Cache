package main

import (
	"context"
	"log"
	"net"
	cacheProto "repMemCache/cache/proto"
	pp "repMemCache/master/proto"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MasterServer struct {
	pp.UnimplementedPeerServiceServer
	cacheProto.UnimplementedCacheServiceServer
	//Peers map[string]*pp.Peer // Changed to map to store peers
	data  map[string]string // Data storage map
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
	log.Print("Hello")
	s.mu.Lock()
	s.data[req.Key] = req.Value

	defer s.mu.Unlock()
	if len(s.Peers.GetPeers()) > 0 {
		for _, peer := range s.Peers.GetPeers() {
			log.Printf("Peers %v", peer)
			conn, err := grpc.Dial(peer.Address, grpc.WithInsecure())
			if err != nil {
				log.Printf("Error in connecting to the peer %v", err)
				continue
			}
			defer conn.Close()
			cacheServiceClient := cacheProto.NewCacheServiceClient(conn) // o peerin cacheine koyuyo valueyu
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

	} else {
		log.Printf("There is no peer in the peerList")
	}

	return &cacheProto.DataResponse{Value: req.Value}, nil
}

func (s *MasterServer) GetData(ctx context.Context, req *cacheProto.KeyRequest) (*cacheProto.DataResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[req.Key]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key not found")
	}
	log.Printf("GetData in master: %v", req.Key)
	return &cacheProto.DataResponse{Value: val}, nil
}

func main() { //master bu
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	peers := pp.NewPeers()
	ms := &MasterServer{Peers: peers, data: make(map[string]string)}
	pp.RegisterPeerServiceServer(grpcServer, ms)
	cacheProto.RegisterCacheServiceServer(grpcServer, ms)
	log.Printf("Master Server is running on port: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
