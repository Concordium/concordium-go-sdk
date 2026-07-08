package tests_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/send"
)

// --- TokenOperationsPayload serialization fix (regression) ---

func TestTokenOperationsPayloadRoundTrip(t *testing.T) {
	payload := &v2.TokenOperationsPayload{
		TokenId:    v2.TokenId{Value: "PLT_TEST_TOKEN"},
		Operations: v2.RawCBOR{Bytes: bytes.Repeat([]byte{0xAB}, 500)}, // > 255 exercises the u32 length
	}

	raw := payload.Encode()
	require.Equal(t, payload.Size()+1, len(raw.Value), "Size()+type byte must equal encoded length")

	decoded, err := raw.Decode()
	require.NoError(t, err)
	tu, ok := decoded.Payload.(v2.TokenUpdate)
	require.True(t, ok)
	require.Equal(t, payload.TokenId, tu.Payload.TokenId)
	require.Equal(t, payload.Operations.Bytes, tu.Payload.Operations.Bytes)
}

// TestPreAccountTransactionTokenUpdateRoundTrip is the exact case that previously
// failed (Encode/Decode width mismatch) and forced the signer node to hand-roll
// header parsing instead of using Deserialize.
func TestPreAccountTransactionTokenUpdateRoundTrip(t *testing.T) {
	sender, err := v2.AccountAddressFromBytes(bytes.Repeat([]byte{0xDE}, 32))
	require.NoError(t, err)

	ops := v2.TokenOperations{
		&v2.BurnOperation{Details: v2.TokenSupplyUpdateDetails{Amount: v2.TokenAmount{Value: 1000, Decimals: 6}}},
	}
	pre, err := construct.TokenUpdateOperation(
		1, sender, v2.SequenceNumber{Value: 1}, v2.TransactionTime{Value: 4_000_000_000},
		v2.TokenId{Value: "PLT_TEST"}, ops,
	)
	require.NoError(t, err)

	ser := pre.Serialize()
	re := v2.PreAccountTransaction{Encoded: &v2.RawPayload{}}
	require.NoError(t, re.Deserialize(ser))
	require.Equal(t, pre.HashToSign.Value, re.HashToSign.Value)

	tu, ok := re.Payload.Payload.(v2.TokenUpdate)
	require.True(t, ok)
	require.Equal(t, "PLT_TEST", tu.Payload.TokenId.Value)
}

// --- UpdateCredentialKeys ---

func TestUpdateCredentialKeysPayloadRoundTrip(t *testing.T) {
	keys := v2.NewCredentialPublicKeysEd25519([][]byte{
		bytes.Repeat([]byte{0x02}, 32),
		bytes.Repeat([]byte{0x03}, 32),
		bytes.Repeat([]byte{0x04}, 32),
	}, 2)
	payload := &v2.UpdateCredentialKeysPayload{
		CredID: v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)},
		Keys:   keys,
	}

	raw := payload.Encode()
	require.Equal(t, payload.Size()+1, len(raw.Value))

	decoded, err := raw.Decode()
	require.NoError(t, err)
	uck, ok := decoded.Payload.(v2.UpdateCredentialKeys)
	require.True(t, ok)
	require.Equal(t, payload.CredID, uck.Payload.CredID)
	require.Equal(t, payload.Keys.Threshold, uck.Payload.Keys.Threshold)
	require.Equal(t, payload.Keys.Keys, uck.Payload.Keys.Keys)
}

// TestUpdateCredentialKeysGoldenBytes pins the exact wire serialization.
func TestUpdateCredentialKeysGoldenBytes(t *testing.T) {
	payload := &v2.UpdateCredentialKeysPayload{
		CredID: v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)},
		Keys: v2.CredentialPublicKeys{
			Keys: map[v2.KeyIndex]v2.VerifyKey{
				0: {Scheme: v2.SchemeEd25519, Key: bytes.Repeat([]byte{0x02}, 32)},
			},
			Threshold: v2.SignatureThreshold{Value: 1},
		},
	}

	var want []byte
	want = append(want, byte(v2.UpdateCredentialKeysPayloadType)) // type tag = 13
	want = append(want, bytes.Repeat([]byte{0x01}, 48)...)        // credId
	want = append(want, 0x01)                                     // numKeys
	want = append(want, 0x00)                                     // keyIndex 0
	want = append(want, 0x00)                                     // scheme Ed25519
	want = append(want, bytes.Repeat([]byte{0x02}, 32)...)        // verify key
	want = append(want, 0x01)                                     // threshold

	require.Equal(t, want, payload.Encode().Value)
}

// TestUpdateCredentialKeysEncodesKeysAscending verifies keys are serialized in
// ascending index order regardless of map iteration order.
func TestUpdateCredentialKeysEncodesKeysAscending(t *testing.T) {
	payload := &v2.UpdateCredentialKeysPayload{
		CredID: v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)},
		Keys: v2.CredentialPublicKeys{
			Keys: map[v2.KeyIndex]v2.VerifyKey{
				2: {Scheme: v2.SchemeEd25519, Key: bytes.Repeat([]byte{0xCC}, 32)},
				0: {Scheme: v2.SchemeEd25519, Key: bytes.Repeat([]byte{0xAA}, 32)},
				1: {Scheme: v2.SchemeEd25519, Key: bytes.Repeat([]byte{0xBB}, 32)},
			},
			Threshold: v2.SignatureThreshold{Value: 2},
		},
	}

	v := payload.Encode().Value
	off := 1 + 48 + 1        // type + credId + numKeys
	const entry = 1 + 1 + 32 // index + scheme + key

	require.Equal(t, byte(0), v[off])
	require.Equal(t, byte(0xAA), v[off+2])
	require.Equal(t, byte(1), v[off+entry])
	require.Equal(t, byte(0xBB), v[off+entry+2])
	require.Equal(t, byte(2), v[off+2*entry])
	require.Equal(t, byte(0xCC), v[off+2*entry+2])
}

func TestConstructUpdateCredentialKeys(t *testing.T) {
	sender, err := v2.AccountAddressFromBytes(bytes.Repeat([]byte{0xDE}, 32))
	require.NoError(t, err)

	credID := v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)}
	keys := v2.NewCredentialPublicKeysEd25519([][]byte{
		bytes.Repeat([]byte{0x02}, 32),
		bytes.Repeat([]byte{0x03}, 32),
		bytes.Repeat([]byte{0x04}, 32),
	}, 2)

	pre, err := construct.UpdateCredentialKeys(
		2, sender, v2.SequenceNumber{Value: 7}, v2.TransactionTime{Value: 4_000_000_000},
		credID, keys, 1,
	)
	require.NoError(t, err)
	require.NotNil(t, pre.HashToSign)

	ser := pre.Serialize()
	re := v2.PreAccountTransaction{Encoded: &v2.RawPayload{}}
	require.NoError(t, re.Deserialize(ser))
	require.Equal(t, pre.HashToSign.Value, re.HashToSign.Value)

	uck, ok := re.Payload.Payload.(v2.UpdateCredentialKeys)
	require.True(t, ok)
	require.Equal(t, credID, uck.Payload.CredID)
	require.Equal(t, uint8(2), uck.Payload.Keys.Threshold.Value)
	require.Len(t, uck.Payload.Keys.Keys, 3)
}

func TestSendUpdateCredentialKeys(t *testing.T) {
	sender, err := v2.AccountAddressFromBytes(bytes.Repeat([]byte{0xDE}, 32))
	require.NoError(t, err)
	keyPair, err := v2.NewKeyPairFromSignKey(bytes.Repeat([]byte{0x11}, 32))
	require.NoError(t, err)
	signer := v2.NewWalletAccount(sender, *keyPair)

	credID := v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)}
	keys := v2.NewCredentialPublicKeysEd25519([][]byte{
		bytes.Repeat([]byte{0x02}, 32),
		bytes.Repeat([]byte{0x03}, 32),
		bytes.Repeat([]byte{0x04}, 32),
	}, 2)

	tx, err := send.UpdateCredentialKeys(
		signer, sender, v2.SequenceNumber{Value: 7}, v2.TransactionTime{Value: 4_000_000_000},
		credID, keys, 1,
	)
	require.NoError(t, err)
	require.Len(t, tx.Signature.Signatures, 1)
	require.Len(t, tx.Signature.Signatures[0].Signatures, 1)

	uck, ok := tx.Payload.Payload.(*v2.UpdateCredentialKeys)
	require.True(t, ok)
	require.Equal(t, credID, uck.Payload.CredID)
	require.Equal(t, uint8(2), uck.Payload.Keys.Threshold.Value)
	require.Len(t, uck.Payload.Keys.Keys, 3)
}

func TestConstructUpdateCredentialKeysValidation(t *testing.T) {
	sender, err := v2.AccountAddressFromBytes(bytes.Repeat([]byte{0xDE}, 32))
	require.NoError(t, err)
	nonce := v2.SequenceNumber{Value: 1}
	exp := v2.TransactionTime{Value: 4_000_000_000}
	credID := v2.CredentialRegistrationID{Value: bytes.Repeat([]byte{0x01}, 48)}
	oneKey := v2.NewCredentialPublicKeysEd25519([][]byte{bytes.Repeat([]byte{0x02}, 32)}, 1)

	// bad credId length
	_, err = construct.UpdateCredentialKeys(1, sender, nonce, exp,
		v2.CredentialRegistrationID{Value: []byte{1, 2, 3}}, oneKey, 1)
	require.Error(t, err)

	// threshold greater than number of keys
	_, err = construct.UpdateCredentialKeys(1, sender, nonce, exp, credID,
		v2.NewCredentialPublicKeysEd25519([][]byte{bytes.Repeat([]byte{0x02}, 32)}, 2), 1)
	require.Error(t, err)

	// non-32-byte key
	badKey := v2.CredentialPublicKeys{
		Keys:      map[v2.KeyIndex]v2.VerifyKey{0: {Scheme: v2.SchemeEd25519, Key: []byte{1, 2, 3}}},
		Threshold: v2.SignatureThreshold{Value: 1},
	}
	_, err = construct.UpdateCredentialKeys(1, sender, nonce, exp, credID, badKey, 1)
	require.Error(t, err)

	// unsupported key scheme
	badScheme := v2.CredentialPublicKeys{
		Keys:      map[v2.KeyIndex]v2.VerifyKey{0: {Scheme: 1, Key: bytes.Repeat([]byte{0x02}, 32)}},
		Threshold: v2.SignatureThreshold{Value: 1},
	}
	_, err = construct.UpdateCredentialKeys(1, sender, nonce, exp, credID, badScheme, 1)
	require.Error(t, err)
}
