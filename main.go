package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/KurobaneShin/blockchain/node"
	"github.com/KurobaneShin/blockchain/proto"
)

func main() {
	node := node.NewNode()

	go func() {
		for {

			time.Sleep(2 * time.Second)
			makeTransaction()

		}
	}()

	log.Fatal(node.Serve(":3000"))
}

func makeTransaction() {
	client, err := grpc.NewClient(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version:    "0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), version)
	if err != nil {
		log.Fatal(err)
	}
}
