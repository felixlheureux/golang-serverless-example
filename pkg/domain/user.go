package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"strings"
	"time"
)

type User struct {
	UserID             string    `db:"user_id" json:"user_id"`
	EthereumAddressHex string    `db:"ethereum_address" json:"ethereum_address"`
	Username           string    `db:"username" json:"username"`
	DefaultCharacterID *string   `db:"default_character_id" json:"default_character_id"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

type UserStoreInput struct {
	EthereumAddressHexInput
	Username string
}

func NewUserStoreInput(addressHex string, Username string) UserStoreInput {
	return UserStoreInput{
		EthereumAddressHexInput: NewEthereumAddressHexInput(addressHex),
		Username:                Username,
	}
}

func (input *UserStoreInput) sanitize() {
	input.Username = strings.TrimSpace(input.Username)
}

func (input UserStoreInput) Validate() error {
	if err := input.EthereumAddressHexInput.Validate(); err != nil {
		return err
	}
	input.sanitize()
	return validation.ValidateStruct(&input,
		validation.Field(&input.Username, validation.Length(3, 20)),
	)
}

type UserUpdateInput struct {
	UserID   string
	Username string
}

func NewUserUpdateInput(userID string, username string) UserUpdateInput {
	return UserUpdateInput{
		UserID:   userID,
		Username: username,
	}
}

func (input *UserUpdateInput) sanitize() {
	input.Username = strings.TrimSpace(input.Username)
}

func (input UserUpdateInput) Validate() error {
	input.sanitize()
	return validation.ValidateStruct(&input,
		validation.Field(&input.Username, validation.Length(3, 20)),
	)
}
