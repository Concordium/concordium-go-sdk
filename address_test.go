package concordium

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
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
		w: []byte(`"foo"`),
	}, {
		n: "ContractAddress",
		a: &Address{contract: &ContractAddress{}},
		w: []byte(`{"index":0,"subindex":0}`),
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
		b: []byte(`"foo"`),
	}, {
		n: "ContractAddress",
		w: &Address{contract: &ContractAddress{}},
		b: []byte(`{"index":0,"subindex":0}`),
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
