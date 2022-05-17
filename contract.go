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
type ModuleRef string

func (m *ModuleRef) Serialize() ([]byte, error) {
	b, err := hex.DecodeString(string(*m))
	if err != nil {
		return nil, fmt.Errorf("%T: hex decode: %w", *m, err)
	}
	if len(b) != moduleRefSize {
		return nil, fmt.Errorf("%T expect %d bytes but %d given", *m, moduleRefSize, len(b))
	}
	return b, nil
}

type ContractName string

type InitName string

func NewInitNameFromContract(contractName ContractName) InitName {
	return InitName(initNamePrefix + contractName)
}

func (n *InitName) Serialize() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.BigEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
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
