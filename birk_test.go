package concordium

import "testing"

const testdataBirkParameters = "testdata/birk_parameters.json"

var testBirkParameters = &BirkParameters{
	ElectionDifficulty: 2.5e-2,
	ElectionNonce:      "139b2dfddd24aa32b5260a8f3908331c814f567ae133de57065d6395db0e8cd6",
	Bakers:             []*BakerInfo{testBakerInfo},
}

func TestBirkParameters_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testBirkParameters, testdataBirkParameters)
}

func TestBirkParameters_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BirkParameters{}, testBirkParameters, testdataBirkParameters)
}
