package payloads

import (
	"errors"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

var (
	// InvalidPayloadType indicated that payload type is invalid.
	InvalidPayloadType = errors.New("invalid payload type")
	// InvalidEncodedPayloadSize indicated that encoded payload size is invalid.
	InvalidEncodedPayloadSize = errors.New("invalid encoded payload size")
)

// PayloadLike is a helper interface so that we can use Payload and EncodedPayload in the same place.
type PayloadLike interface {
	// Encode encodes the transaction payload by serializing.
	Encode() EncodedPayload
	isPayloadLike()
}

// Payload describes payload of an account transaction.
type Payload interface {
	PayloadLike
	isPayload()
}

const PayloadTypeSize = 1

// PayloadType defines Payload type byte.
type PayloadType uint8

const (
	// DeployModulePayloadType defines DeployModulePayload type byte.
	DeployModulePayloadType PayloadType = 0
	// InitContractPayloadType defines InitContractPayload type byte.
	InitContractPayloadType PayloadType = 1
	// UpdateContractPayloadType defines UpdateContractPayload type byte.
	UpdateContractPayloadType PayloadType = 2
	// TransferPayloadType defines TransferPayload type byte.
	TransferPayloadType PayloadType = 3
	// RegisterDataPayloadType defines RegisterDataPayload type byte.
	RegisterDataPayloadType PayloadType = 21
	// TransferWithMemoPayloadType defines TransferWithMemoPayload type byte.
	TransferWithMemoPayloadType PayloadType = 22
)

// GetPayloadType returns PayloadType byte from transmitted Payload.
func GetPayloadType(payload Payload) (PayloadType, error) {
	switch payload.(type) {
	case *DeployModulePayload:
		return DeployModulePayloadType, nil
	case *InitContractPayload:
		return InitContractPayloadType, nil
	case *UpdateContractPayload:
		return UpdateContractPayloadType, nil
	case *TransferPayload:
		return TransferPayloadType, nil
	case *RegisterDataPayload:
		return RegisterDataPayloadType, nil
	case *TransferWithMemoPayload:
		return TransferWithMemoPayloadType, nil
	}
	return 0xff, InvalidPayloadType
}

// decode parses specific Payload from bytes.
func decode(payloadBytes []byte) (payload Payload, err error) {
	if len(payloadBytes) <= 1 {
		return nil, InvalidEncodedPayloadSize
	}

	switch PayloadType(payloadBytes[:PayloadTypeSize][0]) {
	case DeployModulePayloadType:
		deployModulePayload := new(DeployModulePayload)
		err = deployModulePayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = deployModulePayload
	case InitContractPayloadType:
		initContractPayload := new(InitContractPayload)
		err = initContractPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = initContractPayload
	case UpdateContractPayloadType:
		updateContractPayload := new(UpdateContractPayload)
		err = updateContractPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = updateContractPayload
	case TransferPayloadType:
		transferPayload := new(TransferPayload)
		err = transferPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = transferPayload
	case RegisterDataPayloadType:
		registerDataPayload := new(RegisterDataPayload)
		err = registerDataPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = registerDataPayload
	case TransferWithMemoPayloadType:
		transferWithMemoPayload := new(TransferWithMemoPayload)
		err = transferWithMemoPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload = transferWithMemoPayload
	}
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// EncodedPayload describes Payload in serialized state.
type EncodedPayload []byte

func (EncodedPayload) isPayloadLike() {}

// Encode encodes the transaction payload by serializing.
func (encodedPayload EncodedPayload) Encode() EncodedPayload {
	return encodedPayload[:]
}

// Size return size of Encoded Payload.
func (encodedPayload EncodedPayload) Size() types.PayloadSize {
	return types.PayloadSize(len(encodedPayload))
}

// Decode decodes the EncodedPayload into a structured Payload.
// This also checks that all data is used, i.e., that there are no remaining trailing bytes.
func (encodedPayload EncodedPayload) Decode() (Payload, error) {
	return decode(encodedPayload[:])
}

// Serialize returns serialized encoded payload.
func (encodedPayload EncodedPayload) Serialize() []byte {
	return encodedPayload[:]
}
