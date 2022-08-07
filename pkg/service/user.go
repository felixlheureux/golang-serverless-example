package service

import (
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/manta-coder/golang-serverless-example/pkg/store"
	"go.uber.org/zap"
)

type UserService interface {
	Get(userID string) (domain.User, error)
	FindByEthereumAddress(ethereumAddressHex string) (domain.User, error)
	Store(input domain.UserStoreInput) (domain.User, error)
	Update(input domain.UserUpdateInput) (domain.User, error)
	UpdateDefaultCharacter(userID string, characterID string) (domain.User, error)
	Remove(userID string) error
}

type userService struct {
	logger    *zap.SugaredLogger
	userStore store.UserStore
}

func NewUserService(logger *zap.SugaredLogger, userStore store.UserStore) UserService {
	return &userService{logger, userStore}
}

func (s *userService) Get(userID string) (domain.User, error) {
	user, err := s.userStore.Get(userID)
	if err != nil {
		return domain.User{}, domain.ErrUserGetFailed(err)
	}

	return user, nil
}

func (s *userService) FindByEthereumAddress(ethereumAddressHex string) (domain.User, error) {
	user, err := s.userStore.FindByEthereumAddress(ethereumAddressHex)
	if err != nil {
		return domain.User{}, domain.ErrUserFindByEthereumAddressFailed(err)
	}

	return user, nil
}

func (s *userService) Store(input domain.UserStoreInput) (domain.User, error) {
	if err := input.Validate(); err != nil {
		return domain.User{}, err
	}

	result, err := s.userStore.Store(domain.User{
		EthereumAddressHex: input.EthereumAddressHexInput.EthereumAddressHex,
		Username:           input.Username,
	})
	if err != nil {
		return domain.User{}, domain.ErrUserStoreFailed(err)
	}

	return result, nil
}

func (s *userService) Update(input domain.UserUpdateInput) (domain.User, error) {
	if err := input.Validate(); err != nil {
		return domain.User{}, err
	}

	user, err := s.userStore.Get(input.UserID)
	if err != nil {
		return domain.User{}, domain.ErrUserFindByEthereumAddressFailed(err)
	}

	user.Username = input.Username

	result, err := s.userStore.Update(user)
	if err != nil {
		return domain.User{}, domain.ErrUserUpdateFailed(err)
	}

	return result, nil
}

func (s *userService) UpdateDefaultCharacter(userID string, characterID string) (domain.User, error) {
	user, err := s.userStore.Get(userID)
	if err != nil {
		return domain.User{}, domain.ErrUserFindByEthereumAddressFailed(err)
	}

	user.DefaultCharacterID = &characterID

	result, err := s.userStore.Update(user)
	if err != nil {
		return domain.User{}, domain.ErrUserUpdateFailed(err)
	}

	return result, nil
}

func (s *userService) Remove(userID string) error {
	if err := s.userStore.Remove(userID); err != nil {
		return domain.ErrUserRemoveFailed(err)
	}

	return nil
}
