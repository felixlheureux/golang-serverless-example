package store

import (
	"github.com/childrenofukiyo/odin/pkg/auth"
	"github.com/childrenofukiyo/odin/pkg/domain"
	"github.com/childrenofukiyo/odin/pkg/helpers"
	"github.com/childrenofukiyo/odin/pkg/tester"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func createTestChallengeStore() ChallengeStore {
	return NewChallengeStore(tester.GetLogger(), tester.DB())
}

var testChallengeStore = createTestChallengeStore()

func testChallenge(t *testing.T) domain.Challenge {
	t.Helper()

	return domain.Challenge{
		EthereumAddressHex: tester.GenerateEthereumAddress(t),
		Challenge:          helpers.Rand(auth.ChallengeStringLength),
	}
}

func createTestChallenge(t *testing.T) domain.Challenge {
	t.Helper()

	challenge := testChallenge(t)

	challenge, err := testChallengeStore.Store(challenge)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return challenge
}

func TestChallengeStore_Store(t *testing.T) {
	challenge := testChallenge(t)

	createdChallenge, err := testChallengeStore.Store(challenge)
	require.NoError(t, err)

	tester.AssertEqual(t, challenge, createdChallenge)
}

func TestChallengeStore_Get(t *testing.T) {
	challenge := createTestChallenge(t)

	foundChallenge, err := testChallengeStore.Get(challenge.EthereumAddressHex)
	require.NoError(t, err)

	tester.AssertEqual(t, challenge, foundChallenge)

	// should error if no user matches ID
	_, err = testChallengeStore.Get(ksuid.New().String())
	assert.Error(t, err)
}

func TestChallengeStore_Remove(t *testing.T) {
	challenge := createTestChallenge(t)

	err := testChallengeStore.Remove(challenge.EthereumAddressHex)
	require.NoError(t, err)

	// should error if no user matches ID
	_, err = testChallengeStore.Get(challenge.EthereumAddressHex)
	assert.Error(t, err)
}
