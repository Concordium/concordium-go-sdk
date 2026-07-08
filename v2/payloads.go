package v2

import (
	"encoding/binary"
	"errors"
	"sort"
)

var (
	// ErrInvalidPayloadType indicated that payload type is invalid.
	ErrInvalidPayloadType = errors.New("invalid payload type")
	// ErrInvalidRawPayloadSize indicated that raw payload size is invalid.
	ErrInvalidRawPayloadSize = errors.New("invalid raw payload size")
)

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
	// UpdateCredentialKeysPayloadType defines UpdateCredentialKeysPayload type byte.
	UpdateCredentialKeysPayloadType PayloadType = 13
	// RegisterDataPayloadType defines RegisterDataPayload type byte.
	RegisterDataPayloadType PayloadType = 21
	// TransferWithMemoPayloadType defines TransferWithMemoPayload type byte.
	TransferWithMemoPayloadType PayloadType = 22
	// TokenUpdatePayloadType defines TokenUpdatePayload type byte.
	TokenUpdatePayloadType PayloadType = 27
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
	case *UpdateCredentialKeys:
		return UpdateCredentialKeysPayloadType, nil
	case *RegisterData:
		return RegisterDataPayloadType, nil
	case *TransferWithMemo:
		return TransferWithMemoPayloadType, nil
	case *TokenUpdate:
		return TokenUpdatePayloadType, nil
	}
	return 0xff, ErrInvalidPayloadType
}

// decode parses specific Payload from bytes.
func decode(payloadBytes []byte) (payload *AccountTransactionPayload, err error) {
	if len(payloadBytes) <= 1 {
		return nil, ErrInvalidRawPayloadSize
	}

	payload = new(AccountTransactionPayload)
	switch PayloadType(payloadBytes[:PayloadTypeSize][0]) {
	case DeployModulePayloadType:
		deployModulePayload := new(DeployModulePayload)
		err = deployModulePayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = DeployModule{Payload: deployModulePayload}
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
	case UpdateCredentialKeysPayloadType:
		updateCredentialKeysPayload := new(UpdateCredentialKeysPayload)
		err = updateCredentialKeysPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = UpdateCredentialKeys{Payload: updateCredentialKeysPayload}
	case RegisterDataPayloadType:
		registerDataPayload := new(RegisterDataPayload)
		err = registerDataPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = RegisterData{Payload: registerDataPayload}
	case TransferWithMemoPayloadType:
		transferWithMemoPayload := new(TransferWithMemoPayload)
		err = transferWithMemoPayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = TransferWithMemo{Payload: transferWithMemoPayload}
	case TokenUpdatePayloadType:
		tokenUpdatePayload := new(TokenOperationsPayload)
		err = tokenUpdatePayload.Decode(payloadBytes[PayloadTypeSize:])
		payload.Payload = TokenUpdate{Payload: tokenUpdatePayload}

	}
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// RawPayload a pre-serialized payload in the binary serialization format defined by the protocol.
type RawPayload struct {
	Value []byte
}

func (RawPayload) isAccountTransactionPayload() {}

// Encode encodes the transaction payload by serializing.
func (rawPayload RawPayload) Encode() *RawPayload {
	return &rawPayload
}

// Size return size of Encoded Payload.
func (rawPayload RawPayload) Size() PayloadSize {
	return PayloadSize{Value: uint32(len(rawPayload.Value))}
}

// Decode decodes the RawPayload into a structured Payload.
// This also checks that all data is used, i.e., that there are no remaining trailing bytes.
func (rawPayload RawPayload) Decode() (*AccountTransactionPayload, error) {
	if rawPayload.Value == nil {
		return &AccountTransactionPayload{}, ErrInvalidRawPayloadSize
	}
	return decode(rawPayload.Value)
}

// Serialize returns serialized encoded payload.
func (rawPayload RawPayload) Serialize() []byte {
	return rawPayload.Value
}

// DeployModulePayload deploys a Wasm module with the given source.
type DeployModulePayload struct {
	DeployModule *VersionedModuleSource
}

// Encode encodes Payload into RawPayload.
func (payload *DeployModulePayload) Encode() *RawPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(DeployModulePayloadType))
	if payload.DeployModule == nil {
		return &RawPayload{Value: buf}
	}

	module := make([]byte, 0)
	switch m := payload.DeployModule.Module.(type) {
	case ModuleSourceV0:
		buf = append(buf, byte(ModuleVersion0))
		module = m.Value
	case *ModuleSourceV0:
		buf = append(buf, byte(ModuleVersion0))
		module = m.Value
	case ModuleSourceV1:
		buf = append(buf, byte(ModuleVersion1))
		module = m.Value
	case *ModuleSourceV1:
		buf = append(buf, byte(ModuleVersion1))
		module = m.Value
	}
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(module)))
	buf = append(buf, module...)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into DeployModulePayload.
func (payload *DeployModulePayload) Decode(source []byte) error {
	if len(source) <= 5 {
		return ErrInvalidRawPayloadSize
	}

	version := moduleVersion(source[:1][0])
	moduleSize := binary.BigEndian.Uint32(source[1:5])
	if len(source) != int(moduleSize+5) {
		return ErrInvalidRawPayloadSize
	}

	module := make([]byte, 0, moduleSize)
	module = append(module, source[5:]...)
	switch version {
	case ModuleVersion0:
		payload.DeployModule = &VersionedModuleSource{Module: ModuleSourceV0{Value: module}}
	case ModuleVersion1:
		payload.DeployModule = &VersionedModuleSource{Module: ModuleSourceV1{Value: module}}
	default:
		return errors.New("invalid module source version")
	}

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *DeployModulePayload) Size() int {
	// 1 byte (module version) + 4 bytes (source module size) + source module bytes.
	if payload.DeployModule == nil {
		return 0
	}
	return 5 + payload.DeployModule.Size()
}

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

// Encode encodes Payload into RawPayload.
func (payload *InitContractPayload) Encode() *RawPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(InitContractPayloadType))
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	buf = append(buf, payload.ModuleRef.Value[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.InitName.Value)))
	buf = append(buf, payload.InitName.Value...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Parameter.Value)))
	buf = append(buf, payload.Parameter.Value...)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into InitContractPayload.
func (payload *InitContractPayload) Decode(source []byte) error {
	if len(source) <= 42 {
		return ErrInvalidRawPayloadSize
	}

	payload.Amount = &Amount{Value: binary.BigEndian.Uint64(source[:8])}

	if payload.ModuleRef == nil {
		payload.ModuleRef = new(ModuleRef)
	}
	copy(payload.ModuleRef.Value[:], source[8:40])

	initNameSize := binary.BigEndian.Uint16(source[40:42])
	if len(source) <= int(initNameSize+44) {
		return ErrInvalidRawPayloadSize
	}

	payload.InitName = &InitName{Value: string(source[42 : 42+initNameSize])}

	parameterSize := binary.BigEndian.Uint16(source[42+initNameSize : initNameSize+44])
	if len(source) != int(initNameSize+parameterSize+44) {
		return ErrInvalidRawPayloadSize
	}

	payload.Parameter = &Parameter{Value: source[44+initNameSize:]}
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *InitContractPayload) Size() int {
	// 8 bytes (Amount) + 32 bytes (DeployModule reference) + 2 bytes (Init name size) +
	// Init name bytes + 2 bytes (Parameter size) + Parameter bytes.
	if payload.InitName == nil {
		payload.InitName = new(InitName)
	}
	if payload.Parameter == nil {
		payload.Parameter = new(Parameter)
	}

	return 44 + len(payload.InitName.Value) + len(payload.Parameter.Value)
}

// RegisterDataPayload registers the given data on the chain.
type RegisterDataPayload struct {
	// The data to register.
	Data *RegisteredData
}

// Encode encodes Payload into RawPayload.
func (payload *RegisterDataPayload) Encode() *RawPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(RegisterDataPayloadType))
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Data.Value)))
	buf = append(buf, payload.Data.Value...)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into RegisterDataPayload.
func (payload *RegisterDataPayload) Decode(source []byte) error {
	if len(source) <= 2 {
		return ErrInvalidRawPayloadSize
	}

	registerDataSize := binary.BigEndian.Uint16(source[:2])
	if len(source) != int(registerDataSize+2) {
		return ErrInvalidRawPayloadSize
	}

	payload.Data = &RegisteredData{Value: source[2:]}

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *RegisterDataPayload) Size() int {
	// 2 bytes (register data size) + register data bytes.
	if payload.Data == nil {
		payload.Data = new(RegisteredData)
	}

	return 2 + len(payload.Data.Value)
}

// TransferPayload transfers CCD to an account.
type TransferPayload struct {
	// Address to send to.
	Receiver *AccountAddress
	// Amount to send.
	Amount *Amount
}

// Encode encodes Payload into RawPayload.
func (payload *TransferPayload) Encode() *RawPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferPayloadType))
	buf = append(buf, payload.Receiver.Value[:]...)
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into TransferPayload.
func (payload *TransferPayload) Decode(source []byte) error {
	if len(source) != payload.Size() {
		return ErrInvalidRawPayloadSize
	}

	if payload.Receiver == nil {
		payload.Receiver = new(AccountAddress)
	}
	copy(payload.Receiver.Value[:], source[:32])
	payload.Amount = &Amount{Value: binary.BigEndian.Uint64(source[32:])}

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferPayload) Size() int {
	// 32 bytes (account address) + 8 bytes (amount).
	if payload.Receiver == nil {
		payload.Receiver = new(AccountAddress)
	}
	if payload.Amount == nil {
		payload.Amount = new(Amount)
	}

	return 40
}

// TransferWithMemoPayload payload of a transfer between two accounts with a memo.
type TransferWithMemoPayload struct {
	// Address to send to.
	Receiver *AccountAddress
	// Memo to include in the transfer.
	Memo *Memo
	// Amount to send.
	Amount *Amount
}

// Encode encodes Payload into RawPayload.
func (payload *TransferWithMemoPayload) Encode() *RawPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferWithMemoPayloadType))
	buf = append(buf, payload.Receiver.Value[:]...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Memo.Value)))
	buf = append(buf, payload.Memo.Value...)
	buf = binary.BigEndian.AppendUint64(buf, payload.Amount.Value)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into TransferWithMemoPayload.
func (payload *TransferWithMemoPayload) Decode(source []byte) error {
	if len(source) <= 34 {
		return ErrInvalidRawPayloadSize
	}

	if payload.Receiver == nil {
		payload.Receiver = new(AccountAddress)
	}
	copy(payload.Receiver.Value[:], source[:32])

	memoSize := binary.BigEndian.Uint16(source[32:34])
	if len(source) != int(42+memoSize) {
		return ErrInvalidRawPayloadSize
	}

	payload.Memo = &Memo{Value: source[34 : 34+memoSize]}
	payload.Amount = &Amount{Value: binary.BigEndian.Uint64(source[34+memoSize:])}

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferWithMemoPayload) Size() int {
	// 32 bytes (account address) + 2 bytes (memo size) + memo bytes + 8 bytes (amount).
	if payload.Receiver == nil {
		payload.Receiver = new(AccountAddress)
	}
	if payload.Memo == nil {
		payload.Memo = new(Memo)
	}
	if payload.Amount == nil {
		payload.Amount = new(Amount)
	}

	return 42 + len(payload.Memo.Value)
}

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

// Encode encodes Payload into RawPayload.
func (payload *UpdateContractPayload) Encode() *RawPayload {
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
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into UpdateContractPayload.
func (payload *UpdateContractPayload) Decode(source []byte) error {
	if len(source) <= 26 {
		return ErrInvalidRawPayloadSize
	}

	payload.Amount = &Amount{Value: binary.BigEndian.Uint64(source[:8])}

	if payload.Address == nil {
		payload.Address = new(ContractAddress)
	}
	payload.Address.Index = binary.BigEndian.Uint64(source[8:16])
	payload.Address.Subindex = binary.BigEndian.Uint64(source[16:24])

	receiveNameSize := binary.BigEndian.Uint16(source[24:26])
	if len(source) <= int(receiveNameSize+28) {
		return ErrInvalidRawPayloadSize
	}

	payload.ReceiveName = &ReceiveName{Value: string(source[26 : 26+receiveNameSize])}

	parameterSize := binary.BigEndian.Uint16(source[26+receiveNameSize : receiveNameSize+28])
	if len(source) != int(receiveNameSize+parameterSize+28) {
		return ErrInvalidRawPayloadSize
	}

	payload.Parameter = &Parameter{Value: source[28+receiveNameSize:]}
	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *UpdateContractPayload) Size() int {
	// 8 bytes (Amount) + 16 bytes (Contract address) + 2 bytes (Receive name size) +
	// Receive name bytes + 2 bytes (Parameter size) + Parameter bytes.
	if payload.Amount == nil {
		payload.Amount = new(Amount)
	}
	if payload.Address == nil {
		payload.Address = new(ContractAddress)
	}
	if payload.ReceiveName == nil {
		payload.ReceiveName = new(ReceiveName)
	}
	if payload.Parameter == nil {
		payload.Parameter = new(Parameter)
	}

	return 28 + len(payload.ReceiveName.Value) + len(payload.Parameter.Value)
}

// TokenOperationsPayload
// / Payload for protocol level transaction. The transaction is a list of token
// / operations that can be decoded from CBOR using
// / [`TokenOperationsPayload::decode_operations`]. Operations includes
// / governance operations, transfers etc.
type TokenOperationsPayload struct {
	TokenId    TokenId
	Operations RawCBOR
}

type RawCBOR struct {
	Bytes []byte
}

func (payload *TokenOperationsPayload) isAccountTransactionPayload() {}

func (payload *TokenOperationsPayload) Size() int {
	tokenBytes := []byte(payload.TokenId.Value)
	cborBytes := payload.Operations.Bytes

	return 1 + len(tokenBytes) + // TokenID length (u8) + TokenID
		4 + len(cborBytes) // CBOR length (u32) + CBOR
}

func (payload *TokenOperationsPayload) Encode() *RawPayload {
	tokenBytes := []byte(payload.TokenId.Value)
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TokenUpdatePayloadType))
	buf = append(buf, uint8(len(tokenBytes)))
	buf = append(buf, tokenBytes...)
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(payload.Operations.Bytes)))
	buf = append(buf, payload.Operations.Bytes...)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into TokenOperationsPayload. The wire layout mirrors
// Encode: a 1-byte TokenID length, the TokenID, a 4-byte big-endian operations
// length, then the CBOR operations.
func (payload *TokenOperationsPayload) Decode(source []byte) error {
	// Minimum: 1-byte TokenID length + 4-byte operations length.
	if len(source) < 5 {
		return ErrInvalidRawPayloadSize
	}

	// Read TokenID length (u8).
	tokenIDLen := int(source[0])
	if len(source) < 1+tokenIDLen+4 {
		return ErrInvalidRawPayloadSize
	}

	// Read TokenID.
	tokenBytes := source[1 : 1+tokenIDLen]
	payload.TokenId = TokenId{Value: string(tokenBytes)}

	// Read RawCBOR length (u32 big-endian).
	cborStart := 1 + tokenIDLen
	cborLen := int(binary.BigEndian.Uint32(source[cborStart : cborStart+4]))
	if len(source) != cborStart+4+cborLen {
		return ErrInvalidRawPayloadSize
	}

	// Read RawCBOR bytes.
	payload.Operations = RawCBOR{
		Bytes: make([]byte, cborLen),
	}
	copy(payload.Operations.Bytes, source[cborStart+4:])

	return nil
}

const (
	// credentialRegistrationIDSize is the byte length of a CredentialRegistrationID
	// (a compressed BLS12-381 G1 group element).
	credentialRegistrationIDSize = 48
	// ed25519PublicKeySize is the byte length of an Ed25519 public key.
	ed25519PublicKeySize = 32
)

// UpdateCredentialKeysPayload updates the keys and signature threshold of an
// existing credential. Wire layout (after the payload type byte):
//
//	credId    : 48 bytes
//	numKeys   : u8
//	  repeated, ascending key index:
//	    keyIndex  : u8
//	    schemeId  : u8   (Ed25519 = 0)
//	    verifyKey : 32 bytes
//	threshold : u8
type UpdateCredentialKeysPayload struct {
	CredID CredentialRegistrationID
	Keys   CredentialPublicKeys
}

// Size returns the encoded payload size, excluding the payload type byte.
func (payload *UpdateCredentialKeysPayload) Size() int {
	// credId + numKeys(1) + threshold(1) + per-key (index 1 + scheme 1 + key).
	size := len(payload.CredID.Value) + 1 + 1
	for _, vk := range payload.Keys.Keys {
		size += 1 + 1 + len(vk.Key)
	}
	return size
}

// Encode encodes the payload into a RawPayload. Keys are serialized in ascending
// key-index order (Concordium serializes the key map ordered by index).
func (payload *UpdateCredentialKeysPayload) Encode() *RawPayload {
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(UpdateCredentialKeysPayloadType))
	buf = append(buf, payload.CredID.Value...)
	buf = append(buf, uint8(len(payload.Keys.Keys)))

	indices := make([]int, 0, len(payload.Keys.Keys))
	for idx := range payload.Keys.Keys {
		indices = append(indices, int(idx))
	}
	sort.Ints(indices)
	for _, idx := range indices {
		vk := payload.Keys.Keys[KeyIndex(idx)]
		buf = append(buf, byte(idx), byte(vk.Scheme))
		buf = append(buf, vk.Key...)
	}

	buf = append(buf, payload.Keys.Threshold.Value)
	return &RawPayload{Value: buf}
}

// Decode decodes bytes into an UpdateCredentialKeysPayload. Every key is assumed
// to be Ed25519 (32 bytes), the only scheme currently defined.
func (payload *UpdateCredentialKeysPayload) Decode(source []byte) error {
	// credId + numKeys(1) + threshold(1).
	if len(source) < credentialRegistrationIDSize+2 {
		return ErrInvalidRawPayloadSize
	}

	offset := 0
	credID := make([]byte, credentialRegistrationIDSize)
	copy(credID, source[:credentialRegistrationIDSize])
	payload.CredID = CredentialRegistrationID{Value: credID}
	offset += credentialRegistrationIDSize

	numKeys := int(source[offset])
	offset++

	keys := make(map[KeyIndex]VerifyKey, numKeys)
	for i := 0; i < numKeys; i++ {
		// keyIndex(1) + schemeId(1) + key(32).
		if len(source) < offset+2+ed25519PublicKeySize {
			return ErrInvalidRawPayloadSize
		}
		keyIndex := KeyIndex(source[offset])
		scheme := SchemeID(source[offset+1])
		key := make([]byte, ed25519PublicKeySize)
		copy(key, source[offset+2:offset+2+ed25519PublicKeySize])
		keys[keyIndex] = VerifyKey{Scheme: scheme, Key: key}
		offset += 2 + ed25519PublicKeySize
	}

	// Exactly one threshold byte must remain.
	if len(source) != offset+1 {
		return ErrInvalidRawPayloadSize
	}
	payload.Keys = CredentialPublicKeys{
		Keys:      keys,
		Threshold: SignatureThreshold{Value: source[offset]},
	}

	return nil
}
