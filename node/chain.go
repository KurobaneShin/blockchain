package node

import (
	"encoding/hex"

	"github.com/KurobaneShin/blockchain/proto"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (list *HeaderList) Add(h *proto.Header) {
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Height() int {
	return list.Len() - 1
}

func (list *HeaderList) Len() int {
	return len(list.headers)
}

type Chain struct {
	blockStore BlockStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore: bs,
		headers:    NewHeaderList(),
	}
}

func (c *Chain) AddBlock(b *proto.Block) error {
	// add the header to the list of headers
	c.headers.Add(b.Header)
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
