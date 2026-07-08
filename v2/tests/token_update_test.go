package tests_test

import (
	"bytes"
	"encoding/binary"
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

		// Construct raw payload strictly according to TokenOperationsPayload.Decode contract:
		// [1 byte tokenLen][token bytes][4 bytes cborLen][cbor bytes]
		var rawPayload []byte
		rawPayload = append(rawPayload, byte(len(tokenId)))
		rawPayload = append(rawPayload, []byte(tokenId)...)
		rawPayload = binary.BigEndian.AppendUint32(rawPayload, uint32(len(rawCBOR)))
		rawPayload = append(rawPayload, rawCBOR...)

		decodedPayload := &v2.TokenOperationsPayload{}
		err := decodedPayload.Decode(rawPayload)
		require.NoError(t, err)

		// Compare the contents
		require.Equal(t, v2.TokenId{Value: tokenId}, decodedPayload.TokenId)
		require.Equal(t, rawCBOR, decodedPayload.Operations.Bytes)

		// Ensure Encode() of TokenUpdate produces a tagged payload
		encoded := (&v2.TokenUpdate{Payload: decodedPayload}).Encode()
		require.NotNil(t, encoded)
		require.Greater(t, len(encoded.Value), 1)
		// Only assert that the first byte is the TokenUpdate payload type
		require.Equal(t, byte(v2.TokenUpdatePayloadType), encoded.Value[0])
	})

}
