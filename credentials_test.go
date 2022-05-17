package concordium

import (
	"encoding/hex"
	"testing"
)

var (
	testCredSignKey = "b53af4521a678b015bbae217277933e87b978a48a9a07d55cc369cdf5e1ac215"

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

func TestEncryptedSignKey_Decode(t *testing.T) {
	got, err := testCredEncryptedSignKey.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
		return
	}
	if hex.EncodeToString(got) != testCredSignKey {
		t.Errorf("Decode() got = %q, want %q", hex.EncodeToString(got), testCredSignKey)
	}
}
