package concordium

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	testAmount        = NewAmountFromMicroCCD(1)
	testAmountBeBytes = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	testAmountLeBytes = []byte{1, 0, 0, 0, 0, 0, 0, 0}
)

// Equal checks equality of 2 items. Used by github.com/google/go-cmp/cmp package. For tests only!
func (a *Amount) Equal(v *Amount) bool {
	if a != nil && v != nil {
		return a.microCCD == v.microCCD
	}
	return a == nil && v == nil
}

func TestNewAmountZero(t *testing.T) {
	v := uint64(0)
	a := NewAmountZero()
	if a.microCCD != v {
		t.Errorf("NewAmountFromMicroCCD() got = %v, w %v", a.microCCD, v)
	}
}

func TestNewAmountFromMicroCCD(t *testing.T) {
	w := uint64(1)
	a := NewAmountFromMicroCCD(w)
	if a.microCCD != w {
		t.Errorf("NewAmountFromMicroCCD() got = %v, w %v", a.microCCD, w)
	}
}

func TestAmount_MicroCCD(t *testing.T) {
	v := uint64(1)
	a := NewAmountFromMicroCCD(v)
	if a.MicroCCD() != v {
		t.Errorf("MicroCCD() got = %v, w %v", a.MicroCCD(), v)
	}
}

func TestAmount_MarshalJSON(t *testing.T) {
	v := uint64(1)
	w := []byte(fmt.Sprintf(`"%d"`, v))
	a := NewAmountFromMicroCCD(v)
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(b, w) {
		t.Errorf("MarshalJSON() got = %v, w %v", b, w)
	}
}

func TestAmount_UnmarshalJSON(t *testing.T) {
	v := uint64(1)
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

func TestAmount_Serialize(t *testing.T) {
	testSerialize(t, testAmount, testAmountBeBytes)
}

func TestAmount_Deserialize(t *testing.T) {
	testDeserialize(t, &Amount{}, testAmount, testAmountBeBytes)
}

func TestAmount_SerializeModel(t *testing.T) {
	testSerializeModel(t, testAmount, testAmountLeBytes)
}

func TestAmount_DeserializeModel(t *testing.T) {
	testDeserializeModel(t, &Amount{}, testAmount, testAmountLeBytes)
}
