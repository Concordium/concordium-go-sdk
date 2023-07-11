package v2

import (
	"encoding/binary"
	"errors"
)

var (
	// InvalidPayloadType indicated that payload type is invalid.
	InvalidPayloadType = errors.New("invalid payload type")
	// InvalidEncodedPayloadSize indicated that encoded payload size is invalid.
	InvalidEncodedPayloadSize = errors.New("invalid encoded payload size")
)

// Payload describes payload of an account transaction.
type Payload interface {
	// Encode encodes the transaction payload by serializing.
	Encode() EncodedPayload
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

// GetPayloadType returns PayloadType byte from transmitted AccountTransactionPayload.
func GetPayloadType(payload AccountTransactionPayload) (PayloadType, error) {
	switch payload.Payload.(type) {
	case *DeployModule:
		return DeployModulePayloadType, nil
	case *InitContract:
		return InitContractPayloadType, nil
	case *UpdateContract:
		return UpdateContractPayloadType, nil
	case *Transfer:
		return TransferPayloadType, nil
	case *RegisterData:
		return RegisterDataPayloadType, nil
	case *TransferWithMemo:
		return TransferWithMemoPayloadType, nil
	}
	return 0xff, InvalidPayloadType
}

// decode parses specific Payload from bytes.
func decode(payloadBytes []byte) (payload *AccountTransactionPayload, err error) {
	if len(payloadBytes) <= 1 {
		return nil, InvalidEncodedPayloadSize
	}

	switch PayloadType(payloadBytes[:PayloadTypeSize][0]) {
	// TODO: fix it.
	// case DeployModulePayloadType:
	//	deployModulePayload := new(DeployModulePayload)
	//	err = deployModulePayload.Decode(payloadBytes[PayloadTypeSize:])
	//	payload.Payload = deployModulePayload
	case InitContractPayloadType:
		initContractPayload := new(InitContractPayload)
		err = initContractPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = InitContract{Payload: initContractPayload}
	case UpdateContractPayloadType:
		updateContractPayload := new(UpdateContractPayload)
		err = updateContractPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = UpdateContract{Payload: updateContractPayload}
	case TransferPayloadType:
		transferPayload := new(TransferPayload)
		err = transferPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = Transfer{Payload: transferPayload}
	case RegisterDataPayloadType:
		registerDataPayload := new(RegisterDataPayload)
		err = registerDataPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = RegisterData{Payload: registerDataPayload}
	case TransferWithMemoPayloadType:
		transferWithMemoPayload := new(TransferWithMemoPayload)
		err = transferWithMemoPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = TransferWithMemo{Payload: transferWithMemoPayload}
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
func (encodedPayload EncodedPayload) Size() PayloadSize {
	return PayloadSize{Value: uint32(len(encodedPayload))}
}

// Decode decodes the EncodedPayload into a structured Payload.
// This also checks that all data is used, i.e., that there are no remaining trailing bytes.
func (encodedPayload EncodedPayload) Decode() (*AccountTransactionPayload, error) {
	return decode(encodedPayload[:])
}

// Serialize returns serialized encoded payload.
func (encodedPayload EncodedPayload) Serialize() []byte {
	return encodedPayload[:]
}

// Ensure that DeployModulePayload implements Payload.
var _ Payload = (*DeployModulePayload)(nil)

type DeployModulePayload struct {
	Module WasmModule
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

	payload.Module.Version = WasmVersion(source[:1][0])
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

// Ensure that InitContractPayload implements Payload.
var _ Payload = (*InitContractPayload)(nil)

// InitContractPayload contains data needed to initialize a smart contract.
type InitContractPayload struct {
	// Deposit this amount of CCD.
	Amount *Amount
	// Reference to the module from which to initialize the instance.
	ModuleRef *ModuleRef
	// Name of the contract in the module.
	InitName *InitName
	// Parameter to invoke the initialization method with.
	Parameter *Parameter
}

// Encode encodes Payload into EncodedPayload.
func (payload *InitContractPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(InitContractPayloadType))
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	buf = append(buf, payload.ModuleRef.Value[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.InitName.Value)))
	buf = append(buf, payload.InitName.Value...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Parameter.Value)))
	buf = append(buf, payload.Parameter.Value...)
	return buf
}

// Decode decodes bytes into InitContractPayload.
func (payload *InitContractPayload) Decode(source []byte) error {
	if len(source) <= 42 {
		return InvalidEncodedPayloadSize
	}

	payload.Amount.Value = binary.BigEndian.Uint64(source[:8])
	copy(payload.ModuleRef.Value[:], source[8:40])

	initNameSize := binary.BigEndian.Uint16(source[40:42])
	if len(source) <= int(initNameSize+44) {
		return InvalidEncodedPayloadSize
	}

	payload.InitName.Value = string(source[42 : 42+initNameSize])

	parameterSize := binary.BigEndian.Uint16(source[42+initNameSize : initNameSize+44])
	if len(source) != int(initNameSize+parameterSize+44) {
		return InvalidEncodedPayloadSize
	}

	payload.Parameter.Value = source[44+initNameSize:]
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *InitContractPayload) Size() int {
	// 8 bytes (Amount) + 32 bytes (Module reference) + 2 bytes (Init name size) +
	// Init name bytes + 2 bytes (Parameter size) + Parameter bytes.
	return 44 + len(payload.InitName.Value) + len(payload.Parameter.Value)
}

func (*InitContractPayload) isPayload()     {}
func (*InitContractPayload) isPayloadLike() {}

// Ensure that RegisterDataPayload implements Payload.
var _ Payload = (*RegisterDataPayload)(nil)

// RegisterDataPayload registers the given data on the chain.
type RegisterDataPayload struct {
	// The data to register.
	Data *RegisteredData
}

// Encode encodes Payload into EncodedPayload.
func (payload *RegisterDataPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(RegisterDataPayloadType))
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Data.Value)))
	buf = append(buf, payload.Data.Value...)
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

	copy(payload.Data.Value[:], source[2:])

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *RegisterDataPayload) Size() int {
	// 2 bytes (register data size) + register data bytes.
	return 2 + len(payload.Data.Value)
}

func (*RegisterDataPayload) isPayload()     {}
func (*RegisterDataPayload) isPayloadLike() {}

// Ensure that TransferPayload implements Payload.
var _ Payload = (*TransferPayload)(nil)

// TransferPayload transfers CCD to an account.
type TransferPayload struct {
	// Address to send to.
	Receiver *AccountAddress
	// Amount to send.
	Amount *Amount
}

// Encode encodes Payload into EncodedPayload.
func (payload *TransferPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferPayloadType))
	buf = append(buf, payload.Receiver.Value[:]...)
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	return buf
}

// Decode decodes bytes into TransferPayload.
func (payload *TransferPayload) Decode(source []byte) error {
	if len(source) != payload.Size() {
		return InvalidEncodedPayloadSize
	}

	copy(payload.Receiver.Value[:], source[:32])
	payload.Amount.Value = binary.BigEndian.Uint64(source[32:])

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferPayload) Size() int {
	// 32 bytes (account address) + 8 bytes (amount).
	return 40
}

func (*TransferPayload) isPayload()     {}
func (*TransferPayload) isPayloadLike() {}

// Ensure that TransferWithMemoPayload implements Payload.
var _ Payload = (*TransferWithMemoPayload)(nil)

// TransferWithMemoPayload payload of a transfer between two accounts with a memo.
type TransferWithMemoPayload struct {
	// Address to send to.
	Receiver *AccountAddress
	// Memo to include in the transfer.
	Memo *Memo
	// Amount to send.
	Amount *Amount
}

// Encode encodes Payload into EncodedPayload.
func (payload *TransferWithMemoPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferPayloadType))
	buf = append(buf, payload.Receiver.Value[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Memo.Value)))
	buf = append(buf, payload.Memo.Value...)
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	return buf
}

// Decode decodes bytes into TransferWithMemoPayload.
func (payload *TransferWithMemoPayload) Decode(source []byte) error {
	if len(source) <= 34 {
		return InvalidEncodedPayloadSize
	}

	copy(payload.Receiver.Value[:], source[:32])
	memoSize := binary.BigEndian.Uint16(source[32:34])
	if len(source) != int(42+memoSize) {
		return InvalidEncodedPayloadSize
	}

	copy(payload.Memo.Value[:], source[34:34+memoSize])
	payload.Amount.Value = binary.BigEndian.Uint64(source[34+memoSize:])

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferWithMemoPayload) Size() int {
	// 32 bytes (account address) + 2 bytes (memo size) + memo bytes + 8 bytes (amount).
	return 42 + len(payload.Memo.Value)
}

func (*TransferWithMemoPayload) isPayload()     {}
func (*TransferWithMemoPayload) isPayloadLike() {}

// Ensure that UpdateContractPayload implements Payload.
var _ Payload = (*UpdateContractPayload)(nil)

// UpdateContractPayload updates a smart contract instance by invoking a specific function.
type UpdateContractPayload struct {
	// Send the given amount of CCD together with the message to the
	// contract instance.
	Amount *Amount
	// Address of the contract instance to invoke.
	Address *ContractAddress
	// Name of the method to invoke on the contract.
	ReceiveName *ReceiveName
	// Parameter to send to the contract instance.
	Parameter *Parameter
}

// Encode encodes Payload into EncodedPayload.
func (payload *UpdateContractPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(UpdateContractPayloadType))
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	buf = binary.BigEndian.AppendUint64(buf, payload.Address.Index)
	buf = binary.BigEndian.AppendUint64(buf, payload.Address.Subindex)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.ReceiveName.Value)))
	buf = append(buf, payload.ReceiveName.Value...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Parameter.Value)))
	buf = append(buf, payload.Parameter.Value...)
	return buf
}

// Decode decodes bytes into UpdateContractPayload.
func (payload *UpdateContractPayload) Decode(source []byte) error {
	if len(source) <= 26 {
		return InvalidEncodedPayloadSize
	}

	payload.Amount.Value = binary.BigEndian.Uint64(source[:8])
	payload.Address.Index = binary.BigEndian.Uint64(source[8:16])
	payload.Address.Subindex = binary.BigEndian.Uint64(source[16:24])

	receiveNameSize := binary.BigEndian.Uint16(source[24:26])
	if len(source) <= int(receiveNameSize+28) {
		return InvalidEncodedPayloadSize
	}

	payload.ReceiveName.Value = string(source[26 : 26+receiveNameSize])

	parameterSize := binary.BigEndian.Uint16(source[26+receiveNameSize : receiveNameSize+28])
	if len(source) != int(receiveNameSize+parameterSize+28) {
		return InvalidEncodedPayloadSize
	}

	payload.Parameter.Value = source[28+receiveNameSize:]
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *UpdateContractPayload) Size() int {
	// 8 bytes (Amount) + 16 bytes (Contract address) + 2 bytes (Receive name size) +
	// Receive name bytes + 2 bytes (Parameter size) + Parameter bytes.
	return 28 + len(payload.ReceiveName.Value) + len(payload.Parameter.Value)
}

func (*UpdateContractPayload) isPayload()     {}
func (*UpdateContractPayload) isPayloadLike() {}
