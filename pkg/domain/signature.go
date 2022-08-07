package domain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/manta-coder/golang-serverless-example/pkg/helpers"
)

type Signature [SignatureSize]byte

func newSignatureFromBytes(sigBytes []byte) Signature {
	sig := Signature{}
	copy(sig[:], sigBytes[:])

	if sig[SignatureSize-1] < SignatureRIRangeBase {
		sig[SignatureSize-1] += SignatureRIRangeBase
	}

	return sig
}

func NewSignatureFromHex(sigHex string) Signature {
	return newSignatureFromBytes(common.FromHex(sigHex))
}

func (sig Signature) Bytes() []byte {
	return sig[:]
}

func ValidateSignatureHex(sigHex string) error {
	if helpers.HasHexPrefix(sigHex) {
		sigHex = sigHex[2:]
	}

	if len(sigHex) != 2*SignatureSize {
		return ErrInvalidSignatureSize(nil)
	}
	if !helpers.IsHex(sigHex) {
		return ErrInvalidSignatureHex(nil)
	}

	return nil
}
