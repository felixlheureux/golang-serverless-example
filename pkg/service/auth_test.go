package service

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/manta-coder/golang-serverless-example/pkg/auth"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/manta-coder/golang-serverless-example/pkg/store"
	"github.com/manta-coder/golang-serverless-example/pkg/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createTestChallengeStore() store.ChallengeStore {
	return store.NewChallengeStore(tester.GetLogger(), tester.DB())
}

var testChallengeStore = createTestChallengeStore()

const authSecret = "123456789abcdefghijklmnopqrstuvwyz"

func createTestAuth() *auth.Service {
	return auth.NewService(authSecret, time.Duration(900)*time.Second)
}

var testAuth = createTestAuth()

func createTestAuthService() AuthService {
	return NewAuthService(tester.GetLogger(), testAuth, testChallengeStore, testUserService)
}

var testAuthService = createTestAuthService()

func TestAuthService_Challenge(t *testing.T) {
	user := testUser(t)

	createdChallenge, err := testAuthService.Challenge(auth.NewChallengeInput(user.EthereumAddressHex))
	require.NoError(t, err)

	foundChallenge, err := testChallengeStore.Get(user.EthereumAddressHex)
	require.NoError(t, err)

	assert.Equal(t, createdChallenge.Challenge, foundChallenge.Challenge)
}

func TestAuthService_AuthorizeService(t *testing.T) {
	privateKey := tester.CreatePrivateKey(t, "9")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	addressHex := domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA)).Hex()

	createdChallenge, err := testAuthService.Challenge(auth.NewChallengeInput(addressHex))
	require.NoError(t, err)

	signedHash := tester.SignHash(createdChallenge.Challenge)

	signatureBytes, err := crypto.Sign(signedHash.Bytes(), privateKey)
	require.NoError(t, err)

	token, err := testAuthService.Authorize(auth.NewAuthorizeInput(addressHex, hexutil.Encode(signatureBytes)))
	require.NoError(t, err)

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(authSecret), nil
	}

	claims := auth.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token.Token, &claims, keyFunc)
	require.NoError(t, err)

	require.True(t, parsedToken.Valid)

	authClaims := parsedToken.Claims.(*auth.Claims)

	assert.Equal(t, authClaims.EthereumAddressHex, addressHex)

	// wrong signature should fail
	privateKey2 := tester.CreatePrivateKey(t, "8")
	signatureBytes2, err := crypto.Sign(signedHash.Bytes(), privateKey2)
	assert.NoError(t, err)
	_, err = testAuthService.Authorize(auth.NewAuthorizeInput(addressHex, hexutil.Encode(signatureBytes2)))
	assert.Error(t, err)

	// wrong public key should fail
	publicKey2 := privateKey2.Public()
	publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey)
	addressHex2 := domain.EthereumAddress(crypto.PubkeyToAddress(*publicKeyECDSA2)).Hex()
	assert.True(t, ok)
	_, err = testAuthService.Authorize(auth.NewAuthorizeInput(addressHex2, hexutil.Encode(signatureBytes2)))
	assert.Error(t, err)
}
