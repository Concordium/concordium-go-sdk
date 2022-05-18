package concordium

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

const (
	accountAddressVersion = 1

	accountAddressType  = 0
	contractAddressType = 1

	accountAddressSize  = 32
	contractAddressSize = 16

	accountAddressJsonType  = "AddressAccount"
	contractAddressJsonType = "AddressContract"
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
	tmp := struct {
		Type    string `json:"type"`
		Address any    `json:"address"`
	}{}
	if a.account != "" {
		tmp.Type = accountAddressJsonType
		tmp.Address = a.account
		return json.Marshal(tmp)
	}
	if a.contract != nil {
		tmp.Type = contractAddressJsonType
		tmp.Address = a.contract
		return json.Marshal(tmp)
	}
	return nil, nil
}

func (a *Address) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	tmp := &struct {
		Type    string          `json:"type"`
		Address json.RawMessage `json:"address"`
	}{}
	err := json.Unmarshal(b, tmp)
	if err != nil {
		return err
	}
	switch tmp.Type {
	case accountAddressJsonType:
		return json.Unmarshal(tmp.Address, &a.account)
	case contractAddressJsonType:
		a.contract = &ContractAddress{}
		return json.Unmarshal(tmp.Address, a.contract)
	default:
		return fmt.Errorf("%T: unexpected data `%s`", *a, b)
	}
}

func (a *Address) SerializeModel() ([]byte, error) {
	var ser []byte
	var typ uint8
	var err error
	switch true {
	case a.account != "":
		typ = accountAddressType
		ser, err = a.account.SerializeModel()
	case a.contract != nil:
		typ = contractAddressType
		ser, err = a.contract.SerializeModel()
	default:
		err = fmt.Errorf("is empty")
	}
	if err != nil {
		return nil, fmt.Errorf("%T: %w", *a, err)
	}
	b := make([]byte, len(ser)+1)
	b[0] = typ
	copy(b[1:], ser)
	return b, nil
}

func (a *Address) DeserializeModel(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, fmt.Errorf("%T requires %d bytes", *a, 1)
	}
	var i int
	var err error
	switch b[0] {
	case accountAddressType:
		i, err = a.account.DeserializeModel(b[1:])
	case contractAddressType:
		a.contract = &ContractAddress{}
		i, err = a.contract.DeserializeModel(b[1:])
	default:
		err = fmt.Errorf("invalid")
	}
	if err != nil {
		return 0, fmt.Errorf("%T: %w", *a, err)
	}
	return i + 1, nil
}

// AccountAddress base-58 check with version byte 1 encoded address (with Bitcoin mapping table)
type AccountAddress string

func (a *AccountAddress) Serialize() ([]byte, error) {
	b, _, err := base58.CheckDecode(string(*a))
	if err != nil {
		return nil, fmt.Errorf("%T: base58 decode: %w", *a, err)
	}
	if len(b) != accountAddressSize {
		return nil, fmt.Errorf("%T expect %d bytes but %d given", *a, accountAddressSize, len(b))
	}
	return b, nil
}

func (a *AccountAddress) Deserialize(b []byte) error {
	if len(b) < accountAddressSize {
		return fmt.Errorf("%T requires %d bytes", *a, accountAddressSize)
	}
	*a = AccountAddress(base58.CheckEncode(b[:accountAddressSize], accountAddressVersion))
	return nil
}

func (a *AccountAddress) SerializeModel() ([]byte, error) {
	return a.Serialize()
}

func (a *AccountAddress) DeserializeModel(b []byte) (int, error) {
	return accountAddressSize, a.Deserialize(b)
}

// ContractAddress is a JSON record with two fields {index : Int, subindex : Int}
type ContractAddress struct {
	Index    uint64 `json:"index"`
	SubIndex uint64 `json:"subindex"`
}

// Serialize returns bytes (Index and SubIndex in big-endian order) of serialized ContractAddress.
func (a *ContractAddress) Serialize() ([]byte, error) {
	b := make([]byte, contractAddressSize)

	binary.BigEndian.PutUint64(b, a.Index)
	binary.BigEndian.PutUint64(b[8:], a.SubIndex)

	return b, nil
}

func (a *ContractAddress) Deserialize(b []byte) error {
	if len(b) < contractAddressSize {
		return fmt.Errorf("%T requires %d bytes", *a, contractAddressSize)
	}
	a.Index = binary.BigEndian.Uint64(b[:8])
	a.SubIndex = binary.BigEndian.Uint64(b[8:])
	return nil
}

func (a *ContractAddress) SerializeModel() ([]byte, error) {
	b := make([]byte, contractAddressSize)

	binary.LittleEndian.PutUint64(b, a.Index)
	binary.LittleEndian.PutUint64(b[8:], a.SubIndex)

	return b, nil
}

func (a *ContractAddress) DeserializeModel(b []byte) (int, error) {
	if len(b) < contractAddressSize {
		return 0, fmt.Errorf("%T requires %d bytes", *a, contractAddressSize)
	}
	a.Index = binary.LittleEndian.Uint64(b[:8])
	a.SubIndex = binary.LittleEndian.Uint64(b[8:])
	return contractAddressSize, nil
}
