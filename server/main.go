package main

import (
	"context"
	"log"
	"net"
	cache "repMemCache/cache"
	cacheProto "repMemCache/cache/proto"
	pp "repMemCache/master/proto"

	"google.golang.org/grpc"
)

type peerServer struct {
	cacheProto.UnimplementedCacheServiceServer
	Cache *cache.Cache
	pp.UnimplementedPeerServiceServer
}

func (s *peerServer) SetData(ctx context.Context, req *cacheProto.DataRequest) (*cacheProto.DataResponse, error) {
	s.Cache.Set(req.Key, req.Value)
	return &cacheProto.DataResponse{}, nil
}
func (s *peerServer) GetData(ctx context.Context, req *cacheProto.KeyRequest) (*cacheProto.DataResponse, error) {
	value, exists := s.Cache.Get(req.Key)
	if !exists {
		return &cacheProto.DataResponse{}, nil
	}
	return &cacheProto.DataResponse{Value: value}, nil

}

func (s *peerServer) SyncData(ctx context.Context, req *cacheProto.DataRequest) (*cacheProto.DataResponse, error) {
	s.Cache.Set(req.Key, req.Value)
	return &cacheProto.DataResponse{}, nil
}

func (s *peerServer) NotifyPeers(ctx context.Context, req *pp.NotifyPeerRequest) (*pp.NotifyPeerResponse, error) {
	log.Printf("Notify from master: %s", req.Message)
	return &pp.NotifyPeerResponse{Status: "Received"}, nil
}

func registerToMaster(masterAddress, serverAddress string) {
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master server: %v", err)
	}
	defer conn.Close()

	client := pp.NewPeerServiceClient(conn)
	_, err = client.RegisterPeer(context.Background(), &pp.RegisterPeerRequest{Address: serverAddress})
	if err != nil {
		log.Fatalf("Failed to register to master: %v", err)
	}

	log.Println("Registered to master successfully.")
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	newCache := cache.NewCache()
	cacheProto.RegisterCacheServiceServer(grpcServer, &peerServer{Cache: newCache})
	pp.RegisterPeerServiceServer(grpcServer, &peerServer{Cache: newCache})
	registerToMaster("localhost:50052", listener.Addr().String()) //master addr.
	log.Printf("Server is running on port: %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
