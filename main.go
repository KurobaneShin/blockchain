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
	makeNode(":3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"}, false)
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"}, false)

	for {

		time.Sleep(time.Second * 2)
		makeTransaction()
	}

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "1",
		ListenAddr: listenAddr,
	}
	if isValidator {
		cfg.PrivateKey = crypto.GeneratePrivateKey()
	}
	n := node.NewNode(cfg)

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
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
