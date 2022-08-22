package account

import (
	"reflect"
	"testing"
	"time"

	"github.com/Concordium/concordium-go-sdk"
)

type testBody struct{}

func (b *testBody) Serialize() ([]byte, error) {
	return testRandomBytes, nil
}

func (b *testBody) BaseEnergy() int {
	s := headerSize + len(testRandomBytes)
	return 1000 - int(calculateTransactionEnergy(1, s, 0))
}

var (
	testRequest = newRequest(
		testCredentials,
		concordium.MustNewAccountAddressFromString("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63"),
		5,
		time.Unix(1622334455, 0),
		&testBody{},
	)
)

func Test_request_Serialize(t1 *testing.T) {
	var want []byte
	want = append(want, testSignatureBytes...)
	want = append(want, testHeaderBytes...)
	want = append(want, testSignature.bodyBytes...)

	got, err := testRequest.Serialize()
	if err != nil {
		t1.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t1.Errorf("Serialize() got = %v, want %v", got, want)
	}
}
