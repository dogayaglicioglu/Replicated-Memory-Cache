package main

import (
	"context"
	"log"
	"net"
	cache "repMemCache/cache"
	pb "repMemCache/cache/proto"
	pp "repMemCache/protoPeer/proto"

	"google.golang.org/grpc"
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
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master server: %v", err)
	}
	defer conn.Close()

	client := pp.NewPeerServiceClient(conn)
	peers, err := client.ListPeers(context.Background(), &pp.ListPeersRequest{}) //STRING AARAY
	if err != nil {
		log.Fatalf("Failed to register to master: %v", err)
	}
	for _, peer := range peers.Addresses {
		log.Println("Peer:", peer) // Sunucu adreslerini listele
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to peer %s: %v", peer, err)
			continue
		}
		defer conn.Close()

		client := pb.NewCacheServiceClient(conn)
		_, err = client.SyncData(context.Background(), req)
		if err != nil {
			log.Printf("Failed to replicate data to peer %s: %v", peer, err)
		}
	}
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
	listener, err := net.Listen("tcp", ":50059")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	newCache := cache.NewCache()
	peers := pp.NewPeers()
	//peers.AddPeer(listener.Addr().String())
	pb.RegisterCacheServiceServer(grpcServer, &peerServer{Cache: newCache, Peers: peers})
	registerToMaster("localhost:50052", listener.Addr().String()) //master addr.
	log.Printf("Server is running on port: %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
