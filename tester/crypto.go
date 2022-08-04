package tester

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"testing"
)

func GenerateEthereumAddress(t *testing.T) string {
	t.Helper()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}

func CreatePrivateKey(t *testing.T, last string) *ecdsa.PrivateKey {
	t.Helper()

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a1" + last)
	if err != nil {
		t.Fatal(err)
	}

	return privateKey
}

func SignHash(challenge string) common.Hash {
	msg := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(challenge)) + challenge)
	return crypto.Keccak256Hash(msg)
}
