package concordium

import (
	"reflect"
	"testing"
)

var (
	testCredSignKey      = DecryptedSignKey("b53af4521a678b015bbae217277933e87b978a48a9a07d55cc369cdf5e1ac215")
	testCredSignKeyBytes = []byte{
		181, 58, 244, 82, 26, 103, 139, 1, 91, 186, 226, 23, 39, 121, 51, 232,
		123, 151, 138, 72, 169, 160, 125, 85, 204, 54, 156, 223, 94, 26, 194, 21,
	}

	testCredEncryptedSignKey = EncryptedSignKey{
		Metadata: EncryptedSignKeyMetadata{
			Iterations:           100000,
			Salt:                 "QsY4+h31LMs974pPN6QfsA==",
			InitializationVector: "kzyQ24xum3WibCKfvngMlg==",
		},
		Password:   "111111",
		CipherText: "9hTfvFaDb/AYD9xXZ2LVnJ2FrHQhP+daUOP3l6m1tKdP6sPrpvucnA1xcuSgjiX3jfLWCJYEvUMv8oubObe410tJU/PfRZeQeB4xUDs04eE=",
	}
)

func TestDecryptedSignKey_Decode(t *testing.T) {
	got, err := testCredSignKey.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testCredSignKeyBytes) {
		t.Errorf("Decode() got = %v, want %v", got, testCredSignKeyBytes)
	}
}

func TestEncryptedSignKey_Decode(t *testing.T) {
	got, err := testCredEncryptedSignKey.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testCredSignKeyBytes) {
		t.Errorf("Decode() got = %v, want %v", got, testCredSignKeyBytes)
	}
}
