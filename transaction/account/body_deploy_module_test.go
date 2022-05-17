package account

import (
	"reflect"
	"testing"
)

var (
	testModuleDeployBody = newDeployModuleBody(
		testRandomBytes,
	)
	testModuleDeployBodyBytes = []byte{
		0,          // type
		0, 0, 0, 0, // version
		0, 0, 0, 10, // wasm len
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, // wasm body
	}
)

func Test_moduleDeployBody_Serialize(t *testing.T) {
	got, err := testModuleDeployBody.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testModuleDeployBodyBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testModuleDeployBodyBytes)
	}
}
