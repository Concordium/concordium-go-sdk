package concordium

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"
)

func Test_DeserializeModel_Bool(t *testing.T) {
	var v bool
	testDeserializeModelFunc(t, []byte{1}, &v, true)
	testDeserializeModelFunc(t, []byte{0}, &v, false)
}

func Test_DeserializeModel_Int8(t *testing.T) {
	var v int8
	w := int8(1)
	b := []byte{byte(w)}
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Int16(t *testing.T) {
	var v int16
	w := int16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(w))
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Int32(t *testing.T) {
	var v int32
	w := int32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(w))
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Int64(t *testing.T) {
	var v int64
	w := int64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(w))
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Uint8(t *testing.T) {
	var v uint8
	w := uint8(1)
	b := []byte{w}
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Uint16(t *testing.T) {
	var v uint16
	w := uint16(1)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, w)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Uint32(t *testing.T) {
	var v uint32
	w := uint32(1)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, w)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Uint64(t *testing.T) {
	var v uint64
	w := uint64(1)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, w)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_String(t *testing.T) {
	var v string
	w := "foo"
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(w)))
	b = append(b, []byte(w)...)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Slice(t *testing.T) {
	var v []uint8
	w := []uint8{1, 2, 3}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(w)))
	b = append(b, w...)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Map(t *testing.T) {
	var v map[uint8]uint8
	w := map[uint8]uint8{1: 11, 2: 12, 3: 13}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(w)))
	for k, i := range w {
		b = append(b, k, i)
	}
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Struct(t *testing.T) {
	type s struct {
		F1 uint8 `concordium:"model"`
		F2 uint8 `concordium:"model"`
	}
	var v s
	w := s{F1: 10, F2: 11}
	b := []byte{w.F1, w.F2}
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_StructOption(t *testing.T) {
	type s struct {
		F uint8 `concordium:"model,option"`
	}
	var v s
	w1 := s{F: 10}
	b1 := []byte{1, w1.F}
	testDeserializeModelFunc(t, b1, &v, w1)
	w2 := s{}
	b2 := []byte{0}
	testDeserializeModelFunc(t, b2, &v, w2)
}

func Test_DeserializeModel_Custom(t *testing.T) {
	type s struct {
		// check when field contains pointer value
		F1 *AccountAddress `concordium:"model"`
		// check when field contains non-pointer value
		F2 AccountAddress `concordium:"model"`
	}
	var v s
	a := AccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63")
	w := s{
		F1: &a,
		F2: a,
	}
	d, err := a.SerializeModel()
	if err != nil {
		t.Errorf("AccountAddress.SerializeState() error = %v", err)
		return
	}
	var b []byte
	b = append(b, d...)
	b = append(b, d...)
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Time(t *testing.T) {
	var v time.Time
	w := time.Unix(1628075277, 345*1e6)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(w.UnixNano()/1e6))
	testDeserializeModelFunc(t, b, &v, w)
}

func Test_DeserializeModel_Duration(t *testing.T) {
	var v time.Duration
	w := 10 * 24 * time.Hour
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(w/1e6))
	testDeserializeModelFunc(t, b, &v, w)
}

func testDeserializeModelFunc(t *testing.T, b []byte, v, want interface{}) {
	err := DeserializeModel(b, v)
	if err != nil {
		t.Errorf("deserializeState() error = %v", err)
		return
	}
	rv := reflect.ValueOf(v)
	rw := reflect.ValueOf(want)
	if rv.Kind() == reflect.Ptr && rw.Kind() != reflect.Ptr {
		v = rv.Elem().Interface()
	}
	if !reflect.DeepEqual(v, want) {
		t.Errorf("deserializeState() got = %v, want %v", v, want)
	}
}
