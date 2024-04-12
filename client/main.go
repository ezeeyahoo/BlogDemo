package main

import (
	"context"
	"log"

	"github.com/ezeeyahoo/demoBlogServiceInGrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("tcp"+":7878", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := proto.NewBlogServicerClient(conn)

	resp, err := client.GetPost(context.Background(), &proto.GetRequest{PostID: ""})
	if err != nil {
		log.Println("failed to get Post")
	} else {
		log.Println(resp.PostID)
	}

}
