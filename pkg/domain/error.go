package domain

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrCodeUnexpected = NewError(1000, "Internal server error")

	ErrNotFound = NewError(3003, "Not found")

	ErrInvalidEthereumAddressHex = NewError(2001, "ethereum address is not hex")
	ErrInvalidSignatureSize      = NewError(2002, fmt.Sprintf("signature must be %d bytes", SignatureSize))
	ErrInvalidSignatureHex       = NewError(2003, "signature is not hex")
	ErrInvalidSignature          = NewError(2004, "signature is invalid")

	ErrUserGetFailed                    = NewError(3000, "failed to get user")
	ErrUserFindByEthereumAddressFailed  = NewError(3000, "failed to find user by ethereum address user")
	ErrUserStoreFailed                  = NewError(3001, "failed to store user")
	ErrUserUpdateFailed                 = NewError(3002, "failed to update user")
	ErrUserRemoveFailed                 = NewError(3003, "failed to remove user")
	ErrUserUpdateDefaultCharacterFailed = NewError(3003, "failed to update user default character")

	ErrChallengeGetFailed    = NewError(4000, "failed to get challenge")
	ErrChallengeStoreFailed  = NewError(4001, "failed to store challenge")
	ErrChallengeRemoveFailed = NewError(4002, "failed to remove challenge")

	ErrCharactersQueryFailed = NewError(5000, "failed to query characters")
	ErrCharacterClaimFailed  = NewError(5001, "failed to claim character")

	ErrClansQueryFailed = NewError(6000, "failed to query clans")
)

type Error struct {
	Code    int
	Message string
	Cause   error
}

func NewError(code int, msg string) func(error) *Error {
	return func(cause error) *Error {
		e := errors.WithStack(cause)

		return &Error{
			Code:    code,
			Message: msg,
			Cause:   e,
		}
	}
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s  %+v", err.Code, err.Message, err.Cause)
}
