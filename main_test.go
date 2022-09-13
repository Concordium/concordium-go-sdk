package concordium

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var (
	testHexString = "c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	testHexJSON   = []byte(`"` + testHexString + `"`)
	testHex       = MustNewHex(testHexString)
)

func pointer[T any](v T) *T {
	return &v
}

func testTimeMustParse(f, v string) time.Time {
	t, err := time.Parse(f, v)
	if err != nil {
		panic(err)
	}
	return t
}

func testFileUnmarshalJSON(t *testing.T, v, w any, td string) {
	f, err := os.Open(td)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	err = json.NewDecoder(f).Decode(v)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !cmp.Equal(v, w) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", v, w)
	}
}

func testSerialize(t *testing.T, v Serializer, w any) {
	got, err := v.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, w) {
		t.Errorf("Serialize() got = %v, want %v", got, w)
	}
}

func testDeserialize(t *testing.T, v Deserializer, w any, b []byte) {
	err := v.Deserialize(b)
	if err != nil {
		t.Errorf("Deserialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(v, w) {
		t.Errorf("Deserialize() got = %v, want %v", v, w)
	}
}

func testSerializeModel(t *testing.T, v ModelSerializer, w any) {
	got, err := v.SerializeModel()
	if err != nil {
		t.Errorf("SerializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, w) {
		t.Errorf("SerializeModel() got = %v, want %v", got, w)
	}
}

func testDeserializeModel(t *testing.T, v ModelDeserializer, w any, b []byte) {
	_, err := v.DeserializeModel(b)
	if err != nil {
		t.Errorf("DeserializeModel() error = %v", err)
		return
	}
	if !reflect.DeepEqual(v, w) {
		t.Errorf("DeserializeModel() got = %v, want %v", v, w)
	}
}

func TestHex_MarshalJSON(t *testing.T) {
	b, err := json.Marshal(testHex)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(b, testHexJSON) {
		t.Errorf("MarshalJSON() got = %v, w %v", b, testHexJSON)
	}
}

func TestHex_UnmarshalJSON(t *testing.T) {
	var h Hex
	err := json.Unmarshal(testHexJSON, &h)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(h, testHex) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", h, testHex)
	}
}

func TestPairTuple_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		w    any
		v    any
		b    []byte
	}{
		{
			name: "Int and String",
			w: &PairTuple[int, string]{
				First:  10,
				Second: "foo",
			},
			v: &PairTuple[int, string]{},
			b: []byte(`[10, "foo"]`),
		},
		{
			name: "Int and ContractAddress",
			w: &PairTuple[int, *ContractAddress]{
				First: 10,
				Second: &ContractAddress{
					Index:    10,
					SubIndex: 11,
				},
			},
			v: &PairTuple[int, *ContractAddress]{},
			b: []byte(`[10, {"index": 10, "subindex": 11}]`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal(tt.b, tt.v)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if !cmp.Equal(tt.v, tt.w) {
				t.Errorf("UnmarshalJSON() got = %v, w %v", tt.v, tt.w)
			}
		})
	}
}
