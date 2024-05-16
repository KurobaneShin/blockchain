package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	pubKey := privKey.Public()

	assert.Equal(t, pubKeyLen, len(pubKey.Bytes()))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "379ae8f011612048a7849ddff3d2cab7900e0897c0b233feb186161b133176f6"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "ed3daa81f054c3c41bb5c578d8caf32a8d32d5ad"
	)
	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()

	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	// test with invalid msg
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	invalidPrivateKey := GeneratePrivateKey()
	invalidPublicKey := invalidPrivateKey.Public()
	assert.False(t, sig.Verify(invalidPublicKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
}
