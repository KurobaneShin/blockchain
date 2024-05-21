package node

import (
	"encoding/hex"

	"github.com/KurobaneShin/blockchain/proto"
)

type Chain struct {
	blockStore BlockStorer
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore: bs,
	}
}

func (c *Chain) AddBlock(b *proto.Block) error {
	// validatation
	return c.blockStore.Put(b)
}

func (c *Chain) GeyBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GeyBlockByHeight(height int) (*proto.Block, error) {
	return nil, nil
}
