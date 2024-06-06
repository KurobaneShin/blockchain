package types

import (
	"bytes"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
	pb "google.golang.org/protobuf/proto"

	"github.com/KurobaneShin/blockchain/crypto"
	"github.com/KurobaneShin/blockchain/proto"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) TxHash {
	return TxHash{hash: hash}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}

func VerifyBlock(b *proto.Block) bool {
	if !VerifyRootHash(b) {
		return false
	}

	if len(b.PublicKey) != crypto.PubKeyLen {
		return false
	}
	sig := crypto.SignatureFromBytes(b.Signature)
	pubkey := crypto.PublicKeyFromBytes(b.PublicKey)
	hash := HashBlock(b)

	return sig.Verify(pubkey, hash)
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	hash := HashBlock(b)
	sig := pk.Sign(hash)
	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()
	return sig
}

func VerifyRootHash(b *proto.Block) bool {
	tree, err := GetMerkleTree(b)
	if err != nil {
		return false
	}

	isTreeValid, err := tree.VerifyTree()
	if err != nil {
		return false
	}

	return isTreeValid
}

func GetMerkleTree(b *proto.Block) (*merkletree.MerkleTree, error) {
	list := make([]merkletree.Content, len(b.Transactions))
	for i := 0; i < len(b.Transactions); i++ {
		list[i] = NewTxHash(HashTransaction(b.Transactions[i]))
	}
	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// HashBlock return a SHA256 of the header.
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
