package account

import (
	"github.com/Concordium/concordium-go-sdk"
	"reflect"
	"testing"
)

var (
	testInitContractBody = newInitContractBody(
		concordium.NewAmountZero(),
		"2185f30bd3a4a6d29758799cc7d2a1330a98184acea622135014be858c5ff413",
		"contract_name",
	)
	testInitContractBodyBytes = []byte{
		1,                      // type
		0, 0, 0, 0, 0, 0, 0, 0, // amount
		// module reference
		33, 133, 243, 11, 211, 164, 166, 210, 151, 88, 121, 156, 199, 210, 161, 51,
		10, 152, 24, 74, 206, 166, 34, 19, 80, 20, 190, 133, 140, 95, 244, 19,
		0, 18, // contract name size (with init_ prefix)
		105, 110, 105, 116, 95, 99, 111, 110, 116, 114, 97, 99, 116, 95, 110, 97, 109, 101, // contract name (with init_ prefix)
		0, 0, // params len
	}
)

func Test_initContractBody_Serialize(t *testing.T) {
	got, err := testInitContractBody.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testInitContractBodyBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testInitContractBodyBytes)
	}
}
