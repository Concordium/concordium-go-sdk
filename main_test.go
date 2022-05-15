package concordium

import (
	"context"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	// TODO move to env
	testGrpcTarget      = "35.184.87.228:10003"
	testGrpcToken       = "rpcadmin"
	testBlockHash       = BlockHash("4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
	testBlockHeight     = BlockHeight(88794)
	testAccountAddress  = AccountAddress("3djqZmm3jFEfMHXj4RtuTYLfr7VJ5ZwmVGmNot8sbadxFrA5eW")
	testContractAddress = &ContractAddress{Index: 0, SubIndex: 0}
	testModuleRef       = ModuleRef("85a8a9242518e07617763de99e5c6bdf39d82fa534a8838929a2167655002813")

	testBaseClient BaseClient
)

func TestMain(m *testing.M) {
	var err error
	testBaseClient, err = NewBaseClient(context.Background(), testGrpcTarget, testGrpcToken)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

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
