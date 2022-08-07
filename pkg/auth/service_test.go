package auth

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/manta-coder/golang-serverless-example/pkg/tester"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createTestService() *Service {
	return NewService("123456789abcdefghijklmnopqrstuvwyz", time.Duration(900)*time.Second)
}

func TestService_VerifyChallenge(t *testing.T) {
	t.Parallel()

	service := createTestService()
	privateKey := tester.CreatePrivateKey(t, "9")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok)

	challenge := service.NewChallenge()
	signedHash := tester.SignHash(challenge)

	signatureBytes, err := crypto.Sign(signedHash.Bytes(), privateKey)
	assert.NoError(t, err)

	err = service.VerifyChallenge(domain.Challenge{
		ChallengeID:        "",
		Challenge:          challenge,
		EthereumAddressHex: domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA)).Hex(),
	}, signatureBytes)
	assert.NoError(t, err)

	// wrong challenge should fail
	err = service.VerifyChallenge(domain.Challenge{
		ChallengeID:        "",
		Challenge:          service.NewChallenge(),
		EthereumAddressHex: domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA)).Hex(),
	}, signatureBytes)
	assert.Error(t, err)

	// wrong signature should fail
	privateKey2 := tester.CreatePrivateKey(t, "8")
	signatureBytes2, err := crypto.Sign(signedHash.Bytes(), privateKey2)
	assert.NoError(t, err)
	err = service.VerifyChallenge(domain.Challenge{
		ChallengeID:        "",
		Challenge:          service.NewChallenge(),
		EthereumAddressHex: domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA)).Hex(),
	}, signatureBytes2)
	assert.Error(t, err)

	// wrong public key should fail
	publicKey2 := privateKey2.Public()
	publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey)
	assert.True(t, ok)
	err = service.VerifyChallenge(domain.Challenge{
		ChallengeID:        "",
		Challenge:          service.NewChallenge(),
		EthereumAddressHex: domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA2)).Hex(),
	}, signatureBytes)
	assert.Error(t, err)
}
