package tests_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

func TestWalletAccountKeyPair(t *testing.T) {
	// test key.
	privateKeyHex := "123abc456def789acb012def345abc678def901abc234def567abc890def123def456abc789def012abc345def678abc901def234abc567def890abc123def45"
	privateKey, err := hex.DecodeString(privateKeyHex)
	require.NoError(t, err)

	hash := new(v2.TransactionHash)
	copy(hash.Value[:], bytes.Repeat([]byte{1}, 32))

	expectedSignature := "cb486fd02d4e6220849c9b1ab06dec0e4c9c3e3d953caded70f3e3d4141b72d231633933dea33e781dbb66331a34faa465838f0138f67521b20ca74a86ccea01"

	t.Run("invalid key pair", func(t *testing.T) {
		cases := []struct {
			key1      []byte
			key2      []byte
			isTwoKeys bool
		}{
			{key1: nil, isTwoKeys: false},                                   // nil value.
			{key1: []byte{1, 2, 3}, isTwoKeys: false},                       // too short key.
			{key1: privateKey, isTwoKeys: false},                            // too long key.
			{key1: nil, key2: nil, isTwoKeys: true},                         // two nil values.
			{key1: nil, key2: []byte{1, 2, 3}, isTwoKeys: true},             // nil and too short key.
			{key1: []byte{1, 2, 3}, key2: nil, isTwoKeys: true},             // too short key and nil.
			{key1: []byte{1, 2, 3}, key2: []byte{1, 2, 3}, isTwoKeys: true}, // two too short keys.
			{key1: privateKey, key2: privateKey, isTwoKeys: true},           // two too long keys.
			{key1: privateKey[:32], key2: nil, isTwoKeys: true},             // valid key and nil.
			{key1: privateKey[:32], key2: []byte{1, 2, 3}, isTwoKeys: true}, // valid key and too short key.
			{key1: privateKey[:32], key2: privateKey, isTwoKeys: true},      // valid key and too long key.
		}

		for _, test := range cases {
			var err error
			if !test.isTwoKeys {
				_, err = v2.NewKeyPairFromSignKey(test.key1)
			} else {
				_, err = v2.NewKeyPairFromSignKeyAndVerifyKey(test.key1, test.key2)
			}
			require.Error(t, err)
		}
	})

	t.Run("successful sign", func(t *testing.T) {
		keyPair, err := v2.NewKeyPairFromSignKeyAndVerifyKey(privateKey[:32], privateKey[32:])
		require.NoError(t, err)

		walletAccount := v2.NewWalletAccount(v2.AccountAddress{}, *keyPair)

		signature, err := walletAccount.SignTransactionHash(hash)
		require.NoError(t, err)
		require.EqualValues(t, 1, walletAccount.NumberOfKeys())
		require.Equal(t, expectedSignature, hex.EncodeToString(signature.Signatures[0].Signatures[0].Value))
	})

	t.Run("add key pair", func(t *testing.T) {
		keyPair1, err := v2.NewKeyPairFromSignKey(privateKey[:32])
		require.NoError(t, err)

		keyPair2, err := v2.NewKeyPairFromSignKey(privateKey[32:])
		require.NoError(t, err)

		walletAccount := v2.NewWalletAccount(v2.AccountAddress{}, *keyPair1)
		require.EqualValues(t, walletAccount.Keys.Threshold.Value, 1)
		require.Len(t, walletAccount.Keys.Keys, 1)
		require.EqualValues(t, walletAccount.Keys.Keys[0].Threshold.Value, 1)
		require.Len(t, walletAccount.Keys.Keys[0].Keys, 1)

		signature, err := walletAccount.SignTransactionHash(hash)
		require.NoError(t, err)
		require.Len(t, signature.Signatures, 1)
		require.Len(t, signature.Signatures[0].Signatures, 1)

		walletAccount.AddKeyPair(*keyPair2)
		require.EqualValues(t, walletAccount.Keys.Threshold.Value, 1)
		require.Len(t, walletAccount.Keys.Keys, 1)
		require.EqualValues(t, walletAccount.Keys.Keys[0].Threshold.Value, 2)
		require.Len(t, walletAccount.Keys.Keys[0].Keys, 2)

		signature, err = walletAccount.SignTransactionHash(hash)
		require.NoError(t, err)
		require.Len(t, signature.Signatures, 1)
		require.Len(t, signature.Signatures[0].Signatures, 2)
	})
}

func TestReadFromFile(t *testing.T) {
	hash := new(v2.TransactionHash)
	copy(hash.Value[:], bytes.Repeat([]byte{1}, 32))

	expectedSignature := "cb486fd02d4e6220849c9b1ab06dec0e4c9c3e3d953caded70f3e3d4141b72d231633933dea33e781dbb66331a34faa465838f0138f67521b20ca74a86ccea01"

	pathToFile := "./2xBpaHottqhwFZURMZW4uZduQvpxNDSy46iXMYs9kceNGaPpZX.export"
	walletAccount, err := v2.NewWalletAccountFromFile(pathToFile)
	require.NoError(t, err)
	require.Equal(t, "2xBpaHottqhwFZURMZW4uZduQvpxNDSy46iXMYs9kceNGaPpZX", walletAccount.Address.ToBase58())
	require.EqualValues(t, 1, walletAccount.Keys.Threshold.Value)
	require.Len(t, walletAccount.Keys.Keys, 1)
	require.EqualValues(t, 1, walletAccount.Keys.Keys[0].Threshold.Value)
	require.Len(t, walletAccount.Keys.Keys[0].Keys, 1)
	require.Equal(t, "123abc456def789acb012def345abc678def901abc234def567abc890def123d", hex.EncodeToString(walletAccount.Keys.Keys[0].Keys[0].Secret()))
	require.Equal(t, "ef456abc789def012abc345def678abc901def234abc567def890abc123def45", hex.EncodeToString(walletAccount.Keys.Keys[0].Keys[0].Public()))

	signature, err := walletAccount.SignTransactionHash(hash)
	require.NoError(t, err)
	require.EqualValues(t, 1, walletAccount.NumberOfKeys())
	require.Equal(t, expectedSignature, hex.EncodeToString(signature.Signatures[0].Signatures[0].Value))
}
