package concordium

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	account  AccountAddress
	contract *ContractAddress
}

func WrapAccountAddress(addr AccountAddress) *Address {
	return &Address{
		account: addr,
	}
}

func WrapContractAddress(addr *ContractAddress) *Address {
	return &Address{
		contract: addr,
	}
}

func (a *Address) MarshalJSON() ([]byte, error) {
	if a.account != "" {
		return json.Marshal(a.account)
	}
	if a.contract != nil {
		return json.Marshal(a.contract)
	}
	return nil, nil
}

func (a *Address) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	switch b[0] {
	case '"':
		return json.Unmarshal(b, &a.account)
	case '{':
		a.contract = &ContractAddress{}
		return json.Unmarshal(b, a.contract)
	default:
		return fmt.Errorf("%T: unexpected data `%s`", *a, b)
	}
}

// AccountAddress base-58 check with version byte 1 encoded address (with Bitcoin mapping table)
type AccountAddress string

// ContractAddress is a JSON record with two fields {index : Int, subindex : Int}
type ContractAddress struct {
	Index    uint64 `json:"index"`
	SubIndex uint64 `json:"subindex"`
}
