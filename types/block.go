package types

import (
	"crypto/sha256"

	pb "google.golang.org/protobuf/proto"

	"github.com/KurobaneShin/blockchain/crypto"
	"github.com/KurobaneShin/blockchain/proto"
)

func VerifyBlock(b *proto.Block) bool {
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
