package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/KurobaneShin/blockchain/node"
	"github.com/KurobaneShin/blockchain/proto"
)

func main() {

	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {

			time.Sleep(2 * time.Second)
			makeTransaction()

		}
	}()

	proto.RegisterNodeServer(grpcServer, node)
	fmt.Println("node running on port: ", ":3000")

	grpcServer.Serve(ln)
}

func makeTransaction() {
	client, err := grpc.NewClient(":3000", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	_, err = c.HandleTransaction(context.TODO(), &proto.Transaction{})
	if err != nil {
		log.Fatal(err)
	}
}
