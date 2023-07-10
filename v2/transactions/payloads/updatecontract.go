package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that UpdateContractPayload implements Payload.
var _ Payload = (*UpdateContractPayload)(nil)

// UpdateContractPayload updates a smart contract instance by invoking a specific function.
type UpdateContractPayload struct {
	// Send the given amount of CCD together with the message to the
	// contract instance.
	Amount types.Amount
	// Address of the contract instance to invoke.
	Address types.ContractAddress
	// Name of the method to invoke on the contract.
	ReceiveName types.OwnedReceiveName
	// Message to send to the contract instance.
	Message types.OwnedParameter
}

// Encode encodes Payload into EncodedPayload.
func (payload *UpdateContractPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(UpdateContractPayloadType))
	buf = binary.BigEndian.AppendUint64(buf, uint64(payload.Amount))
	buf = append(buf, payload.Address[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.ReceiveName)))
	buf = append(buf, payload.ReceiveName...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Message)))
	buf = append(buf, payload.Message...)
	return buf
}

// Decode decodes bytes into UpdateContractPayload.
func (payload *UpdateContractPayload) Decode(source []byte) error {
	if len(source) <= 26 {
		return InvalidEncodedPayloadSize
	}

	payload.Amount = types.Amount(binary.BigEndian.Uint64(source[:8]))
	copy(payload.Address[:], source[8:24])

	receiveNameSize := binary.BigEndian.Uint16(source[24:26])
	if len(source) <= int(receiveNameSize+28) {
		return InvalidEncodedPayloadSize
	}

	payload.ReceiveName = types.OwnedReceiveName(source[26 : 26+receiveNameSize])

	parameterSize := binary.BigEndian.Uint16(source[26+receiveNameSize : receiveNameSize+28])
	if len(source) != int(receiveNameSize+parameterSize+28) {
		return InvalidEncodedPayloadSize
	}

	payload.Message = source[28+receiveNameSize:]
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *UpdateContractPayload) Size() int {
	// 8 bytes (Amount) + 16 bytes (Contract address) + 2 bytes (Receive name size) +
	// Receive name bytes + 2 bytes (Parameter size) + Parameter bytes.
	return 28 + len(payload.ReceiveName) + len(payload.Message)
}

func (*UpdateContractPayload) isPayload()     {}
func (*UpdateContractPayload) isPayloadLike() {}
