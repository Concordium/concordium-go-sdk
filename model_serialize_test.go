package concordium

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"
)

func Test_SerializeModel_Nil(t *testing.T) {
	testSerializeModelFunc(t, nil, []byte{})
}

func Test_SerializeModel_Bool(t *testing.T) {
	testSerializeModelFunc(t, true, []byte{1})
	testSerializeModelFunc(t, false, []byte{0})
}

func Test_SerializeModel_Int8(t *testing.T) {
	testSerializeModelFunc(t, int8(1), []byte{1})
}

func Test_SerializeModel_Int16(t *testing.T) {
	v := int16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(v))
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Int32(t *testing.T) {
	v := int32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Int64(t *testing.T) {
	v := int64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Uint8(t *testing.T) {
	testSerializeModelFunc(t, uint8(1), []byte{1})
}

func Test_SerializeModel_Uint16(t *testing.T) {
	v := uint16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Uint32(t *testing.T) {
	v := uint32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Uint64(t *testing.T) {
	v := uint64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_String(t *testing.T) {
	v := "foo"
	c := len(v)
	b := make([]byte, 4+c)
	binary.LittleEndian.PutUint32(b, uint32(c))
	copy(b[4:], v)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Slice(t *testing.T) {
	v := []uint8{1, 2, 3}
	c := len(v)
	b := make([]byte, 4+c)
	binary.LittleEndian.PutUint32(b, uint32(c))
	copy(b[4:], v)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Map(t *testing.T) {
	v := map[uint8]uint8{0: 11, 1: 12, 2: 13}
	c := len(v)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c))
	// map iteration order is not guarantied
	for i := 0; i < len(v); i++ {
		b = append(b, uint8(i), v[uint8(i)])
	}
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Struct(t *testing.T) {
	type s struct {
		F1 uint8 `concordium:"model"`
		F2 uint8 `concordium:"model"`
	}
	v := s{F1: 10, F2: 11}
	b := []byte{v.F1, v.F2}
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_StructOption(t *testing.T) {
	type s struct {
		F1 *uint8          `concordium:"model,option"`
		F2 *AccountAddress `concordium:"model,option"`
		F3 *AccountAddress `concordium:"model,option"`
	}
	a := AccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63")
	d, err := a.SerializeModel()
	if err != nil {
		t.Errorf("AccountAddress.ModelSerializer() error = %v", err)
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
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Custom(t *testing.T) {
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
		t.Errorf("AccountAddress.ModelSerializer() error = %v", err)
		return
	}
	var b []byte
	b = append(b, d...)
	b = append(b, d...)
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Time(t *testing.T) {
	v := time.Unix(1628075277, 345*1e6)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v.UnixNano()/1e6))
	testSerializeModelFunc(t, v, b)
}

func Test_SerializeModel_Duration(t *testing.T) {
	v := 10 * 24 * time.Hour
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v/1e6))
	testSerializeModelFunc(t, v, b)
}

func testSerializeModelFunc(t *testing.T, v any, want []byte) {
	b, err := SerializeModel(v)
	if err != nil {
		t.Errorf("SerializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(b, want) {
		t.Errorf("SerializeModel() got = %v, want %v", b, want)
	}
}
