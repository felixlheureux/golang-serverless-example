package auth

import (
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
)

type ChallengeInput struct {
	domain.EthereumAddressHexInput
}

func NewChallengeInput(addressHex string) ChallengeInput {
	return ChallengeInput{
		EthereumAddressHexInput: domain.NewEthereumAddressHexInput(addressHex),
	}
}

type AuthorizeInput struct {
	domain.EthereumAddressHexInput
	SigHex string
}

func NewAuthorizeInput(addressHex, sigHex string) AuthorizeInput {
	return AuthorizeInput{
		EthereumAddressHexInput: domain.NewEthereumAddressHexInput(addressHex),
		SigHex:                  sigHex,
	}
}

func (input AuthorizeInput) Validate() error {
	if err := input.EthereumAddressHexInput.Validate(); err != nil {
		return err
	}
	if err := domain.ValidateSignatureHex(input.SigHex); err != nil {
		return err
	}
	return nil
}

func (input AuthorizeInput) Signature() domain.Signature {
	return domain.NewSignatureFromHex(input.SigHex)
}

type AuthorizeSilentlyInput struct {
	domain.EthereumAddressHexInput
}

func NewAuthorizeSilentlyInput(addressHex string) AuthorizeSilentlyInput {
	return AuthorizeSilentlyInput{
		EthereumAddressHexInput: domain.NewEthereumAddressHexInput(addressHex),
	}
}

func (input AuthorizeSilentlyInput) Validate() error {
	if err := input.EthereumAddressHexInput.Validate(); err != nil {
		return err
	}
	return nil
}
