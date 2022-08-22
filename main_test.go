package concordium

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func testTimeMustParse(f, v string) time.Time {
	t, err := time.Parse(f, v)
	if err != nil {
		panic(err)
	}
	return t
}

func testEqualJSON(a, b []byte) (bool, error) {
	var i, j any
	if err := json.Unmarshal(a, &i); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j, i), nil
}

func testHexMarshalJSON(t *testing.T, v any, w []byte) {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(b, w) {
		t.Errorf("MarshalJSON() got = %v, w %v", b, w)
	}
}

func testHexUnmarshalJSON(t *testing.T, v, w any, b []byte) {
	err := json.Unmarshal(b, v)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(v, w) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", v, w)
	}
}

func testFileMarshalJSON(t *testing.T, v any, td string) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	d, err := os.ReadFile(td)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	ok, err := testEqualJSON(b, d)
	if err != nil {
		t.Fatalf("testEqualJSON() error = %v", err)
	}
	if !ok {
		t.Errorf("MarshalJSON() got = %s, w %s", b, d)
	}
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
