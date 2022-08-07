package store

import (
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/manta-coder/golang-serverless-example/pkg/tester"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createTestUserStore() UserStore {
	return NewUserStore(tester.GetLogger(), tester.DB())
}

var testUserStore = createTestUserStore()

func testUser(t *testing.T) domain.User {
	t.Helper()

	now := time.Now()

	return domain.User{
		EthereumAddressHex: tester.GenerateEthereumAddress(t),
		Username:           ksuid.New().String(),
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

func TestUserStore_Store(t *testing.T) {
	user := testUser(t)

	createdUser, err := testUserStore.Store(user)
	require.NoError(t, err)

	user.UserID = createdUser.UserID
	tester.AssertEqual(t, user, createdUser)
}

func TestUserStore_Get(t *testing.T) {
	user := createTestUser(t)

	foundUser, err := testUserStore.Get(user.UserID)
	require.NoError(t, err)

	tester.AssertEqual(t, user, foundUser)

	// should error if no user matches ID
	_, err = testUserStore.Get(ksuid.New().String())
	assert.Error(t, err)
}

func TestUserStore_FindByEthereumAddress(t *testing.T) {
	user := createTestUser(t)

	foundUser, err := testUserStore.FindByEthereumAddress(user.EthereumAddressHex)
	require.NoError(t, err)

	tester.AssertEqual(t, user, foundUser)
}

func TestUserStore_Update(t *testing.T) {
	user := createTestUser(t)

	updateUser := testUser(t)
	updateUser.UserID = user.UserID

	_, err := testUserStore.Update(updateUser)
	require.NoError(t, err)

	foundUser, err := testUserStore.Get(user.UserID)
	require.NoError(t, err)

	updateUser.EthereumAddressHex = user.EthereumAddressHex
	tester.AssertEqual(t, updateUser, foundUser)
}

func TestUserStore_Remove(t *testing.T) {
	user := createTestUser(t)

	err := testUserStore.Remove(user.UserID)
	require.NoError(t, err)

	// should error if no user matches ID
	_, err = testUserStore.Get(user.UserID)
	assert.Error(t, err)
}
