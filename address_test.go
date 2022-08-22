package concordium

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	testAccountAddressString = "3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63"
	testAccountAddressJSON   = []byte(`"` + testAccountAddressString + `"`)
	testAccountAddress       = AccountAddress{
		67, 216, 242, 23, 249, 75, 83, 21, 191, 33, 90, 180, 74, 75, 37, 207,
		77, 10, 155, 248, 93, 73, 251, 134, 119, 71, 132, 152, 76, 101, 116, 217,
	}

	testContractAddress = &ContractAddress{
		Index:    10,
		SubIndex: 11,
	}
	testContractAddressBeBytes = []byte{0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 11}
	testContractAddressLeBytes = []byte{10, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0}
)

// Equal checks equality of 2 items. Used by github.com/google/go-cmp/cmp package. For tests only!
func (a *Address) Equal(v *Address) bool {
	if a != nil && v != nil {
		return reflect.DeepEqual(a.account, v.account) && a.contract.Equal(v.contract)
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
		a: &Address{account: testAccountAddress},
		w: []byte(`{"type":"AddressAccount","address":"` + testAccountAddressString + `"}`),
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
		w: &Address{account: testAccountAddress},
		b: []byte(`{"type":"AddressAccount","address":"` + testAccountAddressString + `"}`),
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
		want:    append([]byte{0}, testAccountAddress[:]...),
	}, {
		name:    "ContractAddress",
		address: &Address{contract: testContractAddress},
		want:    append([]byte{1}, testContractAddressLeBytes...),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSerializeModel(t, tt.address, tt.want)
		})
	}
}

func TestAddress_DeserializeModel(t *testing.T) {
	tests := []struct {
		name string
		want *Address
		data []byte
	}{{
		name: "AccountAddress",
		want: &Address{account: testAccountAddress},
		data: append([]byte{0}, testAccountAddress[:]...),
	}, {
		name: "ContractAddress",
		want: &Address{contract: testContractAddress},
		data: append([]byte{1}, testContractAddressLeBytes...),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDeserializeModel(t, &Address{}, tt.want, tt.data)
		})
	}
}

func TestAccountAddress_MarshalJSON(t *testing.T) {
	got, err := testAccountAddress.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(got, testAccountAddressJSON) {
		t.Errorf("MarshalJSON() got = %s, w %s", got, testAccountAddressJSON)
	}
}

func TestAccountAddress_UnmarshalJSON(t *testing.T) {
	var a AccountAddress
	err := a.UnmarshalJSON(testAccountAddressJSON)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(a, testAccountAddress) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", a, testAccountAddress)
	}
}

func TestAccountAddress_Serialize(t *testing.T) {
	testSerialize(t, &testAccountAddress, testAccountAddress[:])
}

func TestAccountAddress_Deserialize(t *testing.T) {
	a := testAccountAddress
	testDeserialize(t, &a, &testAccountAddress, testAccountAddress[:])
}

func TestAccountAddress_SerializeModel(t *testing.T) {
	testSerializeModel(t, &testAccountAddress, testAccountAddress[:])
}

func TestAccountAddress_DeserializeModel(t *testing.T) {
	a := testAccountAddress
	testDeserializeModel(t, &a, &testAccountAddress, testAccountAddress[:])
}

func TestContractAddress_Serialize(t *testing.T) {
	testSerialize(t, testContractAddress, testContractAddressBeBytes)
}

func TestContractAddress_Deserialize(t *testing.T) {
	testDeserialize(t, &ContractAddress{}, testContractAddress, testContractAddressBeBytes)
}

func TestContractAddress_SerializeModel(t *testing.T) {
	testSerializeModel(t, testContractAddress, testContractAddressLeBytes)
}

func TestContractAddress_DeserializeModel(t *testing.T) {
	testDeserializeModel(t, &ContractAddress{}, testContractAddress, testContractAddressLeBytes)
}
