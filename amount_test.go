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
		return a.microCCD == v.microCCD
	}
	return a == nil && v == nil
}

func TestAmount_NewAmountZero(t *testing.T) {
	v := uint64(0)
	a := NewAmountZero()
	if a.microCCD != v {
		t.Errorf("NewAmountFromMicroCCD() got = %v, w %v", a.microCCD, v)
	}
}

func TestAmount_NewAmountFromMicroCCD(t *testing.T) {
	v := 1
	w := uint64(v)
	a := NewAmountFromMicroCCD(v)
	if a.microCCD != w {
		t.Errorf("NewAmountFromMicroCCD() got = %v, w %v", a.microCCD, w)
	}
}

func TestAmount_NewAmountFromCCD(t *testing.T) {
	v := 1e-6
	w := uint64(v * 1e6)
	a := NewAmountFromCCD(v)
	if a.microCCD != w {
		t.Errorf("NewAmountFromCCD() got = %v, w %v", a.microCCD, w)
	}
}

func TestAmount_MicroCCD(t *testing.T) {
	v := 1
	a := NewAmountFromMicroCCD(v)
	if a.MicroCCD() != v {
		t.Errorf("MicroCCD() got = %v, w %v", a.MicroCCD(), v)
	}
}

func TestAmount_CCD(t *testing.T) {
	v := 1e-6
	a := NewAmountFromCCD(v)
	if a.CCD() != v {
		t.Errorf("CCD() got = %v, w %v", a.CCD(), v)
	}
}

func TestAmount_MarshalJSON(t *testing.T) {
	v := 1
	w := []byte{byte(v)}
	a := NewAmountFromMicroCCD(v)
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
	if a.MicroCCD() != v {
		t.Errorf("UnmarshalJSON() got = %v, w %v", a.MicroCCD(), v)
	}
}
