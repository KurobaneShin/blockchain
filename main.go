package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/KurobaneShin/blockchain/node"
	"github.com/KurobaneShin/blockchain/proto"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()

	go n.Serve(listenAddr)

	if len(bootstrapNodes) > 0 {
		if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
			log.Fatal(err)
		}
	}

	return n
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
