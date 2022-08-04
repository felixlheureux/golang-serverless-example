package service

import (
	"github.com/childrenofukiyo/odin/pkg/domain"
	"github.com/childrenofukiyo/odin/pkg/helpers"
	"github.com/childrenofukiyo/odin/pkg/store"
	"github.com/childrenofukiyo/odin/pkg/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createTestUserStore() store.UserStore {
	return store.NewUserStore(tester.GetLogger(), tester.DB())
}

var testUserStore = createTestUserStore()

func testUser(t *testing.T) domain.User {
	t.Helper()

	now := time.Now()

	return domain.User{
		EthereumAddressHex: tester.GenerateEthereumAddress(t),
		Username:           helpers.Rand(20),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func createTestUser(t *testing.T) domain.User {
	t.Helper()

	user := testUser(t)

	user, err := testUserStore.Store(user)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return user
}

func createTestUserService() UserService {
	return NewUserService(tester.GetLogger(), testUserStore)
}

var testUserService = createTestUserService()

func TestUserService_Store(t *testing.T) {
	user := testUser(t)

	createdUser, err := testUserService.Store(domain.NewUserStoreInput(user.EthereumAddressHex, user.Username))
	require.NoError(t, err)

	user.UserID = createdUser.UserID
	tester.AssertEqual(t, user, createdUser)
}

func TestUserService_Get(t *testing.T) {
	user := createTestUser(t)

	foundUser, err := testUserService.Get(user.UserID)
	require.NoError(t, err)

	tester.AssertEqual(t, user, foundUser)
}

func TestUserService_FindByEthereumAddress(t *testing.T) {
	user := createTestUser(t)

	foundUser, err := testUserService.FindByEthereumAddress(user.EthereumAddressHex)
	require.NoError(t, err)

	tester.AssertEqual(t, user, foundUser)
}

func TestUserService_Update(t *testing.T) {
	user := createTestUser(t)

	updateUser := testUser(t)
	updateUser.UserID = user.UserID

	_, err := testUserService.Update(domain.NewUserUpdateInput(updateUser.UserID, updateUser.Username))
	require.NoError(t, err)

	foundUser, err := testUserStore.Get(user.UserID)
	require.NoError(t, err)

	tester.AssertEqual(t, updateUser, foundUser)
}

func TestUserService_Remove(t *testing.T) {
	user := createTestUser(t)

	err := testUserService.Remove(user.UserID)
	require.NoError(t, err)

	user, err = testUserService.Get(user.UserID)
	assert.Equal(t, user.UserID, "")
}
