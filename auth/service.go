package auth

import (
	"github.com/childrenofukiyo/odin/pkg/domain"
	"github.com/childrenofukiyo/odin/pkg/helpers"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"time"
)

const (
	ChallengeStringLength = 32
)

type Service struct {
	secret              string
	tokenExpiryDuration time.Duration
}

func NewService(secret string, ted time.Duration) *Service {
	return &Service{
		secret:              secret,
		tokenExpiryDuration: ted,
	}
}

func (s *Service) NewChallenge() string {
	return helpers.Rand(ChallengeStringLength)
}

func (s *Service) VerifyChallenge(userChallenge domain.Challenge, responseBytes []byte) error {
	if responseBytes[domain.SignatureSize-1] >= domain.SignatureRIRangeBase {
		responseBytes[domain.SignatureSize-1] -= domain.SignatureRIRangeBase
	}

	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(userChallenge.Challenge)) + userChallenge.Challenge)
	hash := crypto.Keccak256Hash(hashedMessage)

	publicKey, err := crypto.SigToPub(
		hash.Bytes(),
		responseBytes,
	)
	if err != nil {
		return err
	}

	if address := domain.EthereumAddress(crypto.PubkeyToAddress(*publicKey)); address.Hex() != userChallenge.EthereumAddressHex {
		return domain.ErrInvalidSignature(nil)
	}

	return nil
}

func (s *Service) IssueToken(user domain.User) ([]byte, error) {
	return newToken(user.UserID, domain.NewEthereumAddressFromHex(user.EthereumAddressHex), s.tokenExpiryDuration).signedBytes(s.secret)
}
