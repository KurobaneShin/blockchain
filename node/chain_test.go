package node

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KurobaneShin/blockchain/types"
	"github.com/KurobaneShin/blockchain/util"
)

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	block := util.RandomBlock()
	blockHash := types.HashBlock(block)

	assert.Nil(t, chain.AddBlock(block))
	fetchedBlock, err := chain.GeyBlockByHash(blockHash)
	assert.Nil(t, err)
	assert.Equal(t, block, fetchedBlock)
}
