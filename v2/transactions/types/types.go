package types

const SHA256Len int = 32

type HashBytes [SHA256Len]byte

type TransactionSignHash HashBytes

type TransactionSignature struct {
	Signatures map[CredentialIndex]map[KeyIndex]Signature
}

type CredentialIndex uint8

type KeyIndex uint8

type Signature []byte

type (
	AccountAddress  []byte
	Nonce           uint64
	Energy          uint64
	PayloadSize     uint32
	TransactionTime uint64 // seconds since the unix epoch.
)

func (AccountAddress) Bytes() []byte {
	return []byte{}
}

type GivenEnergy interface {
	isGivenEnergy()
}

type AbsoluteEnergy Energy

func (*AbsoluteEnergy) isGivenEnergy() {}

type AddEnergy struct {
	Energy  Energy
	NumSigs uint32
}

func (*AddEnergy) isGivenEnergy() {}

// An Amount of microCCD.
type Amount uint64

type WasmModule struct {
	Version WasmVersion
	Source  ModuleSource
}

type ModuleSource []byte

func (moduleSource ModuleSource) Size() uint64 {
	return uint64(len(moduleSource))
}

type WasmVersion uint8

const WasmVersion0 WasmVersion = 0
const WasmVersion1 WasmVersion = 1

type ModuleReference [32]byte
type OwnedContractName string
type OwnedParameter []byte

type ContractAddress [16]byte
type OwnedReceiveName string

type RegisterData []byte

type Memo []byte

type CredentialType interface {
	isCredentialType()
}
type CredentialTypeInitial struct{}

func (CredentialTypeInitial) isCredentialType() {}

type CredentialTypeNormal struct{}

func (CredentialTypeNormal) isCredentialType() {}
