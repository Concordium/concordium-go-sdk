package account

import (
	"reflect"
	"testing"
	"time"

	"github.com/Concordium/concordium-go-sdk"
)

var (
	testHeader = &header{
		accountAddress: concordium.MustNewAccountAddress("3TdFQK6hqoqoW38JQJGZ2y3RZfgVPzwB7dMKXbBdeYvdwPeF63"),
		nonce:          5,
		energy:         1000,
		bodySize:       10,
		expiry:         time.Unix(1622334455, 0),
	}
	testHeaderBytes = []byte{
		// address
		67, 216, 242, 23, 249, 75, 83, 21, 191, 33, 90, 180, 74, 75, 37, 207,
		77, 10, 155, 248, 93, 73, 251, 134, 119, 71, 132, 152, 76, 101, 116, 217,
		0, 0, 0, 0, 0, 0, 0, 5, // nonce
		0, 0, 0, 0, 0, 0, 3, 232, // energy
		0, 0, 0, 10, // body len
		0, 0, 0, 0, 96, 178, 219, 247, // expired at
	}
)

func Test_header_Serialize(t *testing.T) {
	got, err := testHeader.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testHeaderBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testHeaderBytes)
	}
}
