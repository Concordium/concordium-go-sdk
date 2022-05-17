package account

import (
	"reflect"
	"testing"
)

var (
	testSignature = &signature{
		cred:        testCredentials,
		headerBytes: testHeaderBytes,
		bodyBytes:   testRandomBytes,
	}
	testSignatureBytes = []byte{
		1,     // credentials count
		0,     // credential index
		1,     // key-pair count
		0,     // key-pair index
		0, 64, // signature len
		// signature
		221, 11, 215, 119, 80, 190, 86, 6, 134, 216, 10, 214, 47, 182, 102, 4, 137, 203, 60, 237, 74, 6, 214, 236, 230, 29, 7, 106, 44, 50, 254, 31,
		134, 143, 114, 148, 150, 134, 45, 9, 124, 27, 125, 214, 75, 255, 194, 26, 91, 216, 207, 217, 80, 239, 55, 61, 58, 190, 180, 109, 95, 55, 212, 2,
	}
)

func Test_signature_Serialize(t *testing.T) {
	got, err := testSignature.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testSignatureBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testSignatureBytes)
	}
}
