package concordium

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
	"time"
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

func testMarshalJSON(t *testing.T, v any, td string) {
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

func testUnmarshalJSON(t *testing.T, v, w any, td string) {
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
