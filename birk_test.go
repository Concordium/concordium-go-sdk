package concordium

import (
	"reflect"
	"testing"
)

const testdataBirkParameters = "testdata/birk_parameters.json"

var (
	testElectionNonceString = "139b2dfddd24aa32b5260a8f3908331c814f567ae133de57065d6395db0e8cd6"
	testElectionNonceJSON   = []byte(`"` + testElectionNonceString + `"`)
	testElectionNonce       = MustNewElectionNonce(testElectionNonceString)

	testBirkParameters = &BirkParameters{
		ElectionDifficulty: 2.5e-2,
		ElectionNonce:      testElectionNonce,
		Bakers:             []*BakerInfo{testBakerInfo},
	}
)

func TestBirkParameters_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BirkParameters{}, testBirkParameters, testdataBirkParameters)
}

func TestElectionNonce_MarshalJSON(t *testing.T) {
	got, err := testElectionNonce.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(got, testElectionNonceJSON) {
		t.Errorf("MarshalJSON() got = %s, w %s", got, testElectionNonceJSON)
	}
}

func TestElectionNonce_UnmarshalJSON(t *testing.T) {
	var n ElectionNonce
	err := n.UnmarshalJSON(testElectionNonceJSON)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(n, testElectionNonce) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", n, testElectionNonce)
	}
}
