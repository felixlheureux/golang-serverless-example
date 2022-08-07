package service

import (
	"github.com/manta-coder/golang-serverless-example/pkg/auth"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/manta-coder/golang-serverless-example/pkg/store"
	"go.uber.org/zap"
)

type AuthService interface {
	Challenge(input auth.ChallengeInput) (auth.ChallengeOutput, error)
	Authorize(input auth.AuthorizeInput) (auth.AuthorizeOutput, error)
	AuthorizeSilently(input auth.AuthorizeSilentlyInput) (auth.AuthorizeOutput, error)
}

type authService struct {
	logger         *zap.SugaredLogger
	auth           *auth.Service
	challengeStore store.ChallengeStore
	userService    UserService
}

func NewAuthService(logger *zap.SugaredLogger, auth *auth.Service, challengeStore store.ChallengeStore, userService UserService) AuthService {
	return &authService{logger, auth, challengeStore, userService}
}

func (s *authService) Challenge(input auth.ChallengeInput) (auth.ChallengeOutput, error) {
	if err := input.Validate(); err != nil {
		return auth.ChallengeOutput{}, err
	}

	address := input.Address()
	challenge := s.auth.NewChallenge()

	if _, err := s.challengeStore.Store(domain.Challenge{
		EthereumAddressHex: address.Hex(),
		Challenge:          challenge,
	}); err != nil {
		return auth.ChallengeOutput{}, domain.ErrChallengeStoreFailed(err)
	}

	return auth.NewChallengeOutput(challenge), nil
}

func (s *authService) Authorize(input auth.AuthorizeInput) (auth.AuthorizeOutput, error) {
	if err := input.Validate(); err != nil {
		return auth.AuthorizeOutput{}, err
	}

	address := input.Address()
	sig := input.Signature()

	challenge, err := s.challengeStore.Get(address.Hex())
	if err != nil {
		return auth.AuthorizeOutput{}, domain.ErrChallengeGetFailed(err)
	}

	verifyErr := s.auth.VerifyChallenge(challenge, sig.Bytes())
	if err = s.challengeStore.Remove(address.Hex()); err != nil {
		return auth.AuthorizeOutput{}, domain.ErrChallengeRemoveFailed(err)
	}
	if verifyErr != nil {
		return auth.AuthorizeOutput{}, verifyErr
	}

	user, err := s.userService.FindByEthereumAddress(address.Hex())
	if err != nil {
		return auth.AuthorizeOutput{}, domain.ErrUserFindByEthereumAddressFailed(err)
	}

	if user.UserID == "" {
		user, err = s.userService.Store(domain.NewUserStoreInput(address.Hex(), ""))
		if err != nil {
			return auth.AuthorizeOutput{}, domain.ErrUserStoreFailed(err)
		}
	}

	tokenBytes, err := s.auth.IssueToken(user)
	if err != nil {
		return auth.AuthorizeOutput{}, err
	}

	return auth.NewAuthorizeOutput(string(tokenBytes)), nil
}

func (s *authService) AuthorizeSilently(input auth.AuthorizeSilentlyInput) (auth.AuthorizeOutput, error) {
	if err := input.Validate(); err != nil {
		return auth.AuthorizeOutput{}, err
	}

	user, err := s.userService.FindByEthereumAddress(input.EthereumAddressHex)
	if err != nil {
		return auth.AuthorizeOutput{}, domain.ErrUserFindByEthereumAddressFailed(err)
	}

	tokenBytes, err := s.auth.IssueToken(user)
	if err != nil {
		return auth.AuthorizeOutput{}, err
	}

	return auth.NewAuthorizeOutput(string(tokenBytes)), nil
}
