package account

import (
	"reflect"
	"testing"

	"github.com/Concordium/concordium-go-sdk"
)

var (
	testSimpleTransferBody = newSimpleTransferBody(
		concordium.MustNewAccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63"),
		concordium.NewAmountZero(),
	)
	testSimpleTransferBodyBytes = []byte{
		3, // type
		// address
		67, 216, 242, 23, 249, 75, 83, 21, 191, 33, 90, 180, 74, 75, 37, 207,
		77, 10, 155, 248, 93, 73, 251, 134, 119, 71, 132, 152, 76, 101, 116, 217,
		0, 0, 0, 0, 0, 0, 0, 0, // amount
	}
)

func Test_simpleTransferBody_Serialize(t *testing.T) {
	got, err := testSimpleTransferBody.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testSimpleTransferBodyBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testSimpleTransferBodyBytes)
	}
}
