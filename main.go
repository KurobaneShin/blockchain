package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/KurobaneShin/blockchain/crypto"
	"github.com/KurobaneShin/blockchain/node"
	"github.com/KurobaneShin/blockchain/proto"
	"github.com/KurobaneShin/blockchain/util"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"})

	time.Sleep(time.Second)
	makeTransaction()

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()

	go n.Serve(listenAddr, bootstrapNodes)

	return n
}

func makeTransaction() {
	client, err := grpc.NewClient(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	privKey := crypto.GeneratePrivateKey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
