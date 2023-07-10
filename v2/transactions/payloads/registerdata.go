package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that RegisterDataPayload implements Payload.
var _ Payload = (*RegisterDataPayload)(nil)

// RegisterDataPayload registers the given data on the chain.
type RegisterDataPayload struct {
	// The data to register.
	Data types.RegisterData
}

// Encode encodes Payload into EncodedPayload.
func (payload *RegisterDataPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(RegisterDataPayloadType))
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Data)))
	buf = append(buf, payload.Data...)
	return buf
}

// Decode decodes bytes into RegisterDataPayload.
func (payload *RegisterDataPayload) Decode(source []byte) error {
	if len(source) <= 2 {
		return InvalidEncodedPayloadSize
	}

	registerDataSize := binary.BigEndian.Uint64(source[32:])
	if len(source) != int(registerDataSize+2) {
		return InvalidEncodedPayloadSize
	}

	copy(payload.Data[:], source[2:])

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *RegisterDataPayload) Size() int {
	// 2 bytes (register data size) + register data bytes.
	return 2 + len(payload.Data)
}

func (*RegisterDataPayload) isPayload()     {}
func (*RegisterDataPayload) isPayloadLike() {}
