package types

import (
	"crypto/sha256"

	pb "google.golang.org/protobuf/proto"

	"github.com/KurobaneShin/blockchain/crypto"
	"github.com/KurobaneShin/blockchain/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTranscation(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		if len(input.Signature) == 0 {
			panic("the transaction has no signature")
		}
		var (
			sig    = crypto.SignatureFromBytes(input.Signature)
			pubKey = crypto.PublicKeyFromBytes(input.PublicKey)
		)

		tmpSig := input.Signature
		input.Signature = nil

		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}

		input.Signature = tmpSig
	}

	return true
}
