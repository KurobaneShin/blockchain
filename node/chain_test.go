package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KurobaneShin/blockchain/crypto"
	"github.com/KurobaneShin/blockchain/proto"
	"github.com/KurobaneShin/blockchain/types"
	"github.com/KurobaneShin/blockchain/util"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	require.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)

		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)

		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))
		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestAddBlockWithTxInsuficientFunds(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedString(seed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("4519e2c0cd2c4f14e753affa6b3ea98ceee2b02ddbfff60ab75481e6bb93d563")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}

	outputs := []*proto.TxOutput{
		{
			Amount:  1001,
			Address: recipient,
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.Error(t, chain.AddBlock(block))
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedString(seed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("4519e2c0cd2c4f14e753affa6b3ea98ceee2b02ddbfff60ab75481e6bb93d563")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}

	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}
