package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that InitContractPayload implements Payload.
var _ Payload = (*InitContractPayload)(nil)

// InitContractPayload contains data needed to initialize a smart contract.
type InitContractPayload struct {
	// Deposit this amount of CCD.
	Amount types.Amount
	// Reference to the module from which to initialize the instance.
	ModRef types.ModuleReference
	// Name of the contract in the module.
	InitName types.OwnedContractName
	// Message to invoke the initialization method with.
	Param types.OwnedParameter
}

// Encode encodes Payload into EncodedPayload.
func (payload *InitContractPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(InitContractPayloadType))
	buf = binary.BigEndian.AppendUint64(buf, uint64(payload.Amount))
	buf = append(buf, payload.ModRef[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.InitName)))
	buf = append(buf, payload.InitName...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Param)))
	buf = append(buf, payload.Param...)
	return buf
}

// Decode decodes bytes into InitContractPayload.
func (payload *InitContractPayload) Decode(source []byte) error {
	if len(source) <= 42 {
		return InvalidEncodedPayloadSize
	}

	payload.Amount = types.Amount(binary.BigEndian.Uint64(source[:8]))
	copy(payload.ModRef[:], source[8:40])

	initNameSize := binary.BigEndian.Uint16(source[40:42])
	if len(source) <= int(initNameSize+44) {
		return InvalidEncodedPayloadSize
	}

	payload.InitName = types.OwnedContractName(source[42 : 42+initNameSize])

	parameterSize := binary.BigEndian.Uint16(source[42+initNameSize : initNameSize+44])
	if len(source) != int(initNameSize+parameterSize+44) {
		return InvalidEncodedPayloadSize
	}

	payload.Param = source[44+initNameSize:]
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *InitContractPayload) Size() int {
	// 8 bytes (Amount) + 32 bytes (Module reference) + 2 bytes (Init name size) +
	// Init name bytes + 2 bytes (Parameter size) + Parameter bytes.
	return 44 + len(payload.InitName) + len(payload.Param)
}

func (*InitContractPayload) isPayload()     {}
func (*InitContractPayload) isPayloadLike() {}
