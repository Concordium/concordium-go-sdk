package concordium

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// Equal checks equality of 2 items. Used by github.com/google/go-cmp/cmp package. For tests only!
func (a *Amount) Equal(v *Amount) bool {
	if a != nil && v != nil {
		return a.microGTU == v.microGTU
	}
	return a == nil && v == nil
}

func TestAmount_NewAmountZero(t *testing.T) {
	v := uint64(0)
	a := NewAmountZero()
	if a.microGTU != v {
		t.Errorf("NewAmountFromMicroGTU() got = %v, w %v", a.microGTU, v)
	}
}

func TestAmount_NewAmountFromMicroGTU(t *testing.T) {
	v := 1
	w := uint64(v)
	a := NewAmountFromMicroGTU(v)
	if a.microGTU != w {
		t.Errorf("NewAmountFromMicroGTU() got = %v, w %v", a.microGTU, w)
	}
}

func TestAmount_NewAmountFromGTU(t *testing.T) {
	v := 1e-6
	w := uint64(v * 1e6)
	a := NewAmountFromGTU(v)
	if a.microGTU != w {
		t.Errorf("NewAmountFromGTU() got = %v, w %v", a.microGTU, w)
	}
}

func TestAmount_MicroGTU(t *testing.T) {
	v := 1
	a := NewAmountFromMicroGTU(v)
	if a.MicroGTU() != v {
		t.Errorf("MicroGTU() got = %v, w %v", a.MicroGTU(), v)
	}
}

func TestAmount_GTU(t *testing.T) {
	v := 1e-6
	a := NewAmountFromGTU(v)
	if a.GTU() != v {
		t.Errorf("GTU() got = %v, w %v", a.GTU(), v)
	}
}

func TestAmount_MarshalJSON(t *testing.T) {
	v := 1
	w := []byte{byte(v)}
	a := NewAmountFromMicroGTU(v)
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if reflect.DeepEqual(b, w) {
		t.Errorf("MarshalJSON() got = %s, w %s", b, w)
	}
}

func TestAmount_UnmarshalJSON(t *testing.T) {
	v := 1
	b := []byte(fmt.Sprintf(`"%d"`, v))
	a := &Amount{}
	err := json.Unmarshal(b, a)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if a.MicroGTU() != v {
		t.Errorf("UnmarshalJSON() got = %v, w %v", a.MicroGTU(), v)
	}
}
