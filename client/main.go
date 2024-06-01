package main

import (
	"context"
	"log"
	pb "repMemCache/cache/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	creds := credentials.NewTLS(nil) //nil means insecure
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Error in connect to the client: %v", err)
	}
	defer conn.Close()

	client := pb.NewCacheServiceClient(conn)

	_, err = client.SetData(context.Background(), &pb.DataRequest{Key: "key1", Value: "value1"})
	if err != nil {
		log.Fatalf("Error calling SetData: %v", err)
	}

	res, err := client.GetData(context.Background(), &pb.KeyRequest{Key: "key1"})
	if err != nil {
		log.Fatalf("Error calling GetData: %v", err)
	}
	log.Printf("GetData response: %s", res.Value)
}
