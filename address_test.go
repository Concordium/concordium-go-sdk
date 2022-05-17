package concordium

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

var (
	testAccountAddress      = AccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63")
	testAccountAddressBytes = []byte{
		67, 216, 242, 23, 249, 75, 83, 21, 191, 33, 90, 180, 74, 75, 37, 207,
		77, 10, 155, 248, 93, 73, 251, 134, 119, 71, 132, 152, 76, 101, 116, 217,
	}

	testContractAddress = &ContractAddress{
		Index:    10,
		SubIndex: 11,
	}
	testContractAddressBigEndianBytes = []byte{0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 11}
	testContractAddressLitEndianBytes = []byte{10, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0}
)

// Equal checks equality of 2 items. Used by github.com/google/go-cmp/cmp package. For tests only!
func (a *Address) Equal(v *Address) bool {
	if a != nil && v != nil {
		return a.account == v.account && a.contract.Equal(v.contract)
	}
	return a == nil && v == nil
}

// Equal checks equality of 2 items. Used by github.com/google/go-cmp/cmp package. For tests only!
func (a *ContractAddress) Equal(v *ContractAddress) bool {
	if a != nil && v != nil {
		return a.Index == v.Index && a.SubIndex == v.SubIndex
	}
	return a == nil && v == nil
}

func TestAddress_MarshalJSON(t *testing.T) {
	tests := []struct {
		n string
		a *Address
		w []byte
	}{{
		n: "AccountAddress",
		a: &Address{account: "foo"},
		w: []byte(`{"type":"AddressAccount","address":"foo"}`),
	}, {
		n: "ContractAddress",
		a: &Address{contract: &ContractAddress{}},
		w: []byte(`{"type":"AddressContract","address":{"index":0,"subindex":0}}`),
	}, {
		n: "Empty",
		a: &Address{},
		w: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.n, func(t *testing.T) {
			a := &Address{
				account:  tt.a.account,
				contract: tt.a.contract,
			}
			got, err := a.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.w) {
				t.Errorf("MarshalJSON() got = %s, w %s", got, tt.w)
			}
		})
	}
}

func TestAddress_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		n string
		b []byte
		w *Address
	}{{
		n: "AccountAddress",
		w: &Address{account: "foo"},
		b: []byte(`{"type":"AddressAccount","address":"foo"}`),
	}, {
		n: "ContractAddress",
		w: &Address{contract: &ContractAddress{}},
		b: []byte(`{"type":"AddressContract","address":{"index":0,"subindex":0}}`),
	}, {
		n: "Empty",
		w: &Address{},
		b: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.n, func(t *testing.T) {
			a := &Address{}
			err := a.UnmarshalJSON(tt.b)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if !cmp.Equal(a, tt.w) {
				t.Errorf("UnmarshalJSON() got = %v, w %v", a, tt.w)
			}
		})
	}
}

func TestAddress_SerializeModel(t *testing.T) {
	tests := []struct {
		name    string
		address *Address
		want    []byte
	}{{
		name:    "AccountAddress",
		address: &Address{account: testAccountAddress},
		want:    append([]byte{0}, testAccountAddressBytes...),
	}, {
		name:    "ContractAddress",
		address: &Address{contract: testContractAddress},
		want:    append([]byte{1}, testContractAddressLitEndianBytes...),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.address.SerializeModel()
			if err != nil {
				t.Errorf("SerializeModel() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeModel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountAddress_Serialize(t *testing.T) {
	got, err := testAccountAddress.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testAccountAddressBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testAccountAddressBytes)
	}
}

func TestAccountAddress_SerializeModel(t *testing.T) {
	got, err := testAccountAddress.SerializeModel()
	if err != nil {
		t.Errorf("SerializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testAccountAddressBytes) {
		t.Errorf("SerializeModel() got = %v, want %v", got, testAccountAddressBytes)
	}
}

func TestContractAddress_Serialize(t *testing.T) {
	got, err := testContractAddress.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testContractAddressBigEndianBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testContractAddressBigEndianBytes)
	}
}

func TestContractAddress_SerializeModel(t *testing.T) {
	got, err := testContractAddress.SerializeModel()
	if err != nil {
		t.Errorf("SerializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testContractAddressLitEndianBytes) {
		t.Errorf("SerializeModel() got = %v, want %v", got, testContractAddressLitEndianBytes)
	}
}
