package main

import (
	"context"
	"log"
	cache "repMemCache/cache/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50059", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error in connect to the client: %v", err)
	}
	defer conn.Close()

	client := cache.NewCacheServiceClient(conn)

	_, err = client.SetData(context.Background(), &cache.DataRequest{Key: "key35", Value: "value25"})
	if err != nil {
		log.Fatalf("Error calling SetData: %v", err)
	}

	res, err := client.GetData(context.Background(), &cache.KeyRequest{Key: "key35"})
	if err != nil {
		log.Fatalf("Error calling GetData: %v", err)
	}
	log.Printf("GetData response: %s", res.Value)
}
