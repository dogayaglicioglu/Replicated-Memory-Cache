package main

import (
	"context"
	"log"
	"net"
	cache "repMemCache/cache"
	pb "repMemCache/cache/proto"
	pp "repMemCache/protoPeer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type peerServer struct {
	pb.UnimplementedCacheServiceServer
	Cache *cache.Cache
	Peers *pp.Peers
}

func (s *peerServer) SetData(ctx context.Context, req *pb.DataRequest) (*pb.DataResponse, error) {
	s.Cache.Set(req.Key, req.Value)
	go s.replicateData(&pb.DataRequest{Key: req.Key, Value: req.Value})
	return &pb.DataResponse{}, nil
}
func (s *peerServer) GetData(ctx context.Context, req *pb.KeyRequest) (*pb.DataResponse, error) {
	value, exists := s.Cache.Get(req.Key)
	if !exists {
		return &pb.DataResponse{}, nil
	}
	return &pb.DataResponse{Value: value}, nil

}

func (s *peerServer) SyncData(ctx context.Context, req *pb.DataRequest) (*pb.DataResponse, error) {
	s.Cache.Set(req.Key, req.Value)
	return &pb.DataResponse{}, nil
}

func (s *peerServer) replicateData(req *pb.DataRequest) {
	creds := credentials.NewTLS(nil) //nil means insecure
	for _, peer := range s.Peers.GetPeers() {
		conn, err := grpc.NewClient(peer.Address, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Printf("Failed to connect to peer %s: %v", peer.Address, err)
			continue
		}
		defer conn.Close()

		client := pb.NewCacheServiceClient(conn)
		_, err = client.SyncData(context.Background(), req)
		if err != nil {
			log.Printf("Failed to replicate data to peer %s: %v", peer.Address, err)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	newCache := cache.NewCache()
	peers := pp.NewPeers()
	peers.AddPeer(listener.Addr().String())
	pb.RegisterCacheServiceServer(grpcServer, &peerServer{Cache: newCache, Peers: peers})

	log.Printf("Server is running on port: %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
