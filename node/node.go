package node

import (
	"context"
	"fmt"

	"github.com/KurobaneShin/blockchain/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.None, error) {

	peer,_ := peer.FromContext(ctx)
	fmt.Println("received tx from:",peer)
	return nil, nil
}

func NewNode() *Node {
	return &Node{}
}
