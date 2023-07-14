package tests_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

func TestSigner(t *testing.T) {
	// test key.
	privateKeyHex := "123abc456def789acb012def345abc678def901abc234def567abc890def123def456abc789def012abc345def678abc901def234abc567def890abc123def45"
	privateKey, err := hex.DecodeString(privateKeyHex)
	require.NoError(t, err)

	hash := new(v2.TransactionHash)
	copy(hash.Value[:], bytes.Repeat([]byte{1}, 32))

	expectedSignature := "cb486fd02d4e6220849c9b1ab06dec0e4c9c3e3d953caded70f3e3d4141b72d231633933dea33e781dbb66331a34faa465838f0138f67521b20ca74a86ccea01"

	t.Run("negative sign", func(t *testing.T) {
		invalidKey := privateKey[3:]
		signer1 := v2.NewSimpleSigner(invalidKey)

		_, err := signer1.SignTransactionHash(hash)
		require.Error(t, err)
		require.Equal(t, "invalid private key size", err.Error())

		signer2 := v2.NewSimpleSigner(nil)

		_, err = signer2.SignTransactionHash(hash)
		require.Error(t, err)
		require.Equal(t, "invalid private key size", err.Error())
	})

	t.Run("successful sign", func(t *testing.T) {
		signer := v2.NewSimpleSigner(privateKey)

		signature, err := signer.SignTransactionHash(hash)
		require.NoError(t, err)
		require.EqualValues(t, 1, signer.NumberOfKeys())
		require.Equal(t, expectedSignature, hex.EncodeToString(signature.Signatures[0].Signatures[0].Value))
	})
}
