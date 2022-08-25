package concordium

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	InvokeContractResultSuccess InvokeContractResultTag = "success"
	InvokeContractResultFailure InvokeContractResultTag = "failure"

	moduleRefSize       = 32
	initNamePrefix      = "init_"
	receiveNameSplitter = "."
)

type InvokeContractResultTag string

// ModuleRef base-16 encoded module reference (64 characters)
type ModuleRef [moduleRefSize]byte

func NewModuleRefFromString(s string) (ModuleRef, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ModuleRef{}, fmt.Errorf("hex decode: %w", err)
	}
	if len(b) != moduleRefSize {
		return ModuleRef{}, fmt.Errorf("expect %d bytes but %d given", moduleRefSize, len(b))
	}
	var h ModuleRef
	copy(h[:], b)
	return h, nil
}

func MustNewModuleRefFromString(s string) ModuleRef {
	h, err := NewModuleRefFromString(s)
	if err != nil {
		panic("MustNewModuleRefFromString: " + err.Error())
	}
	return h
}

func (r *ModuleRef) String() string {
	return hex.EncodeToString((*r)[:])
}

func (r ModuleRef) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(r[:])
	if err != nil {
		return nil, fmt.Errorf("%T: %w", r, err)
	}
	return b, nil
}

func (r *ModuleRef) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *r, err)
	}
	var x ModuleRef
	copy(x[:], v)
	*r = x
	return nil
}

func (r *ModuleRef) Serialize() ([]byte, error) {
	return (*r)[:], nil
}

func (r *ModuleRef) Deserialize(b []byte) error {
	var v ModuleRef
	copy(v[:], b)
	*r = v
	return nil
}

type ContractName string

func (n *ContractName) SerializeModel() ([]byte, error) {
	return nil, fmt.Errorf("use %T instead of %T", InitName(""), *n)
}

func (n *ContractName) DeserializeModel([]byte) (int, error) {
	return 0, fmt.Errorf("use %T instead of %T", InitName(""), *n)
}

type InitName string

func NewInitName(contractName ContractName) InitName {
	return InitName(initNamePrefix + contractName)
}

func (n *InitName) Serialize() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.BigEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *InitName) SerializeModel() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.LittleEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *InitName) DeserializeModel(b []byte) (int, error) {
	i := 2
	if len(b) < i {
		return 0, fmt.Errorf("%T requires %d bytes", *n, i)
	}
	l := int(binary.LittleEndian.Uint16(b))
	*n = InitName(b[i : i+l])
	return i + l, nil
}

type ReceiveName string

func NewReceiveName(contractName ContractName, receiver string) ReceiveName {
	return ReceiveName(string(contractName) + receiveNameSplitter + receiver)
}

func (n *ReceiveName) Serialize() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.BigEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *ReceiveName) SerializeModel() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.LittleEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *ReceiveName) DeserializeModel(b []byte) (int, error) {
	i := 2
	if len(b) < i {
		return 0, fmt.Errorf("%T requires %d bytes", *n, i)
	}
	l := int(binary.LittleEndian.Uint16(b))
	*n = ReceiveName(b[i : i+l])
	return i + l, nil
}

type ContractContext struct {
	Invoker   *Address         `json:"invoker"`
	Contract  *ContractAddress `json:"contract"`
	Amount    *Amount          `json:"amount"`
	Method    ReceiveName      `json:"method"`
	Parameter Model            `json:"parameter"`
	Energy    int              `json:"energy"`
}

type InvokeContractResult struct {
	Tag         InvokeContractResultTag `json:"tag"`
	UsedEnergy  int                     `json:"usedEnergy"`
	ReturnValue Model                   `json:"returnValue"`
	Events      any                     `json:"events"` // TODO
	Reason      any                     `json:"reason"` // TODO
}
