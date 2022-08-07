package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"time"
)

type Claims struct {
	UserID             string `json:"user_id"`
	EthereumAddressHex string `json:"ethereum_address"`
	jwt.StandardClaims
}

func newClaims(userID string, address domain.EthereumAddress, d time.Duration) *Claims {
	now := time.Now()

	return &Claims{
		UserID:             userID,
		EthereumAddressHex: address.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(d).Unix(),
			IssuedAt:  now.Unix(),
		},
	}
}

type token struct {
	*jwt.Token
}

func newToken(userID string, address domain.EthereumAddress, d time.Duration) *token {
	return &token{jwt.NewWithClaims(
		jwt.SigningMethodHS256, newClaims(userID, address, d),
	)}
}

func (token *token) signedString(secret string) (string, error) {
	return token.Token.SignedString([]byte(secret))
}

func (token *token) signedBytes(secret string) ([]byte, error) {
	str, err := token.signedString(secret)
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
