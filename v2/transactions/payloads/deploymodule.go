package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that DeployModulePayload implements Payload.
var _ Payload = (*DeployModulePayload)(nil)

type DeployModulePayload struct {
	Module types.WasmModule
}

// Encode encodes Payload into EncodedPayload.
func (payload *DeployModulePayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(DeployModulePayloadType))
	buf = append(buf, byte(payload.Module.Version))
	buf = binary.BigEndian.AppendUint32(buf, uint32(payload.Module.Source.Size()))
	buf = append(buf, payload.Module.Source...)
	return buf
}

// Decode decodes bytes into DeployModulePayload.
func (payload *DeployModulePayload) Decode(source []byte) error {
	if len(source) <= 5 {
		return InvalidEncodedPayloadSize
	}

	payload.Module.Version = types.WasmVersion(source[:1][0])
	sourceModuleSize := binary.BigEndian.Uint32(source[1:5])
	if len(source) != int(sourceModuleSize+5) {
		return InvalidEncodedPayloadSize
	}

	payload.Module.Source = source[5:]
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *DeployModulePayload) Size() int {
	// 1 byte (module version) + 4 bytes (source module size) + source module bytes.
	return 5 + len(payload.Module.Source)
}

func (*DeployModulePayload) isPayload()     {}
func (*DeployModulePayload) isPayloadLike() {}
