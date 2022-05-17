package concordium

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"
)

func Test_serializeModel_Nil(t *testing.T) {
	testSerializeModel(t, nil, []byte{})
}

func Test_serializeModel_Bool(t *testing.T) {
	testSerializeModel(t, true, []byte{1})
	testSerializeModel(t, false, []byte{0})
}

func Test_SerializeModel_Int8(t *testing.T) {
	testSerializeModel(t, int8(1), []byte{1})
}

func Test_SerializeModel_Int16(t *testing.T) {
	v := int16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(v))
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Int32(t *testing.T) {
	v := int32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Int64(t *testing.T) {
	v := int64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Uint8(t *testing.T) {
	testSerializeModel(t, uint8(1), []byte{1})
}

func Test_SerializeModel_Uint16(t *testing.T) {
	v := uint16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Uint32(t *testing.T) {
	v := uint32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Uint64(t *testing.T) {
	v := uint64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_String(t *testing.T) {
	v := "foo"
	c := len(v)
	b := make([]byte, 4+c)
	binary.LittleEndian.PutUint32(b, uint32(c))
	copy(b[4:], v)
	testSerializeModel(t, v, b)
}

func Test_SerializeModel_Slice(t *testing.T) {
	v := []uint8{1, 2, 3}
	c := len(v)
	b := make([]byte, 4+c)
	binary.LittleEndian.PutUint32(b, uint32(c))
	copy(b[4:], v)
	testSerializeModel(t, v, b)
}

func Test_serializeModel_Map(t *testing.T) {
	v := map[uint8]uint8{0: 11, 1: 12, 2: 13}
	c := len(v)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c))
	// map iteration order is not guarantied
	for i := 0; i < len(v); i++ {
		b = append(b, uint8(i), v[uint8(i)])
	}
	testSerializeModel(t, v, b)
}

func Test_serializeModel_Struct(t *testing.T) {
	type s struct {
		F1 uint8 `concordium:"model"`
		F2 uint8 `concordium:"model"`
	}
	v := s{F1: 10, F2: 11}
	b := []byte{v.F1, v.F2}
	testSerializeModel(t, v, b)
}

func Test_serializeModel_StructOption(t *testing.T) {
	type s struct {
		F1 *uint8          `concordium:"model,option"`
		F2 *AccountAddress `concordium:"model,option"`
		F3 *AccountAddress `concordium:"model,option"`
	}
	a := AccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63")
	d, err := a.SerializeModel()
	if err != nil {
		t.Errorf("AccountAddress.SerializeModel() error = %v", err)
		return
	}
	v := s{
		F3: &a,
	}
	var b []byte
	b = append(b, 0)    // F1
	b = append(b, 0)    // F2
	b = append(b, 1)    // F3
	b = append(b, d...) // F3
	testSerializeModel(t, v, b)
}

func Test_serializeModel_Custom(t *testing.T) {
	type s struct {
		// check when field contains pointer value
		F1 *AccountAddress `concordium:"model"`
		// check when field contains non-pointer value
		F2 AccountAddress `concordium:"model"`
	}
	a := AccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63")
	v := s{
		F1: &a,
		F2: a,
	}
	d, err := a.SerializeModel()
	if err != nil {
		t.Errorf("AccountAddress.SerializeModel() error = %v", err)
		return
	}
	var b []byte
	b = append(b, d...)
	b = append(b, d...)
	testSerializeModel(t, v, b)
}

func Test_serializeModel_Time(t *testing.T) {
	v := time.Unix(1628075277, 345*1e6)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v.UnixNano()/1e6))
	testSerializeModel(t, v, b)
}

func Test_serializeModel_Duration(t *testing.T) {
	v := 10 * 24 * time.Hour
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v/1e6))
	testSerializeModel(t, v, b)
}

func testSerializeModel(t *testing.T, v any, want []byte) {
	b, err := serializeModel(v)
	if err != nil {
		t.Errorf("serializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(b, want) {
		t.Errorf("serializeModel() got = %v, want %v", b, want)
	}
}
