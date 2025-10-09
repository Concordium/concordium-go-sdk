package tests_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Concordium/concordium-go-sdk/v2"
)

func TestUpdateTokenPayloads(t *testing.T) {

	moduleRef := new(v2.ModuleRef)
	copy(moduleRef.Value[:], bytes.Repeat([]byte{1, 2}, 16))

	t.Run("updateToken encode/decode", func(t *testing.T) {
		tokenId := "010203"
		rawCBOR := []byte{0xa1, 0x63, 0x66, 0x6f, 0x6f, 0x63, 0x62, 0x61, 0x72} // CBOR map: {"foo": "bar"}

		tokenUpdatePayload := &v2.TokenUpdate{
			Payload: &v2.TokenOperationsPayload{
				TokenId: v2.TokenId{Value: tokenId},
				Operations: v2.RawCBOR{
					Bytes: rawCBOR,
				},
			},
		}

		encoded := tokenUpdatePayload.Encode()
		require.NotNil(t, encoded)
		require.NotEmpty(t, encoded.Value)

		// Remove the leading type byte (TokenUpdatePayloadType) to decode only the payload
		require.True(t, len(encoded.Value) > 1)
		rawPayload := encoded.Value[1:]

		decodedPayload := &v2.TokenOperationsPayload{}
		err := decodedPayload.Decode(rawPayload)
		require.NoError(t, err)

		// Compare the contents
		require.Equal(t, tokenUpdatePayload.Payload.TokenId, decodedPayload.TokenId)
		require.Equal(t, tokenUpdatePayload.Payload.Operations.Bytes, decodedPayload.Operations.Bytes)

		// Check that re-encoding matches
		reEncoded := decodedPayload.Encode()
		require.Equal(t, encoded.Size(), reEncoded.Size())
		require.Equal(t, encoded.Value, reEncoded.Value)
	})

}
