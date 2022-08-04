package domain

import "github.com/ethereum/go-ethereum/common"

type EthereumAddress common.Address

func NewEthereumAddressFromHex(addressHex string) EthereumAddress {
	return EthereumAddress(common.HexToAddress(addressHex))
}

func (address EthereumAddress) Hex() string {
	return common.Address(address).Hex()
}

func ValidateEthereumAddressHex(address string) error {
	if !common.IsHexAddress(address) {
		return ErrInvalidEthereumAddressHex(nil)
	}
	return nil
}

type EthereumAddressHexInput struct {
	EthereumAddressHex string
}

func NewEthereumAddressHexInput(EthereumAddressHex string) EthereumAddressHexInput {
	return EthereumAddressHexInput{
		EthereumAddressHex: EthereumAddressHex,
	}
}

func (input EthereumAddressHexInput) Validate() error {
	if err := ValidateEthereumAddressHex(input.EthereumAddressHex); err != nil {
		return err
	}
	return nil
}

func (input EthereumAddressHexInput) Address() EthereumAddress {
	return NewEthereumAddressFromHex(input.EthereumAddressHex)
}
