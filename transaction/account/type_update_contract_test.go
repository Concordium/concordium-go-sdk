package account

import (
	"github.com/Concordium/concordium-go-sdk"
	"reflect"
	"testing"
)

var (
	testUpdateContractBody = newUpdateContractBody(
		concordium.NewAmountZero(),
		&concordium.ContractAddress{
			Index:    10,
			SubIndex: 11,
		},
		concordium.NewReceiveName("contract_name", "receive_name"),
		uint8(8), uint16(16), uint32(32), uint64(64),
	)
	testUpdateContractBodyBytes = []byte{
		2,                      // type
		0, 0, 0, 0, 0, 0, 0, 0, // amount
		0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 11, // contract address
		0, 26, // receive function name len (contract_name.receive_name)
		// receive function name (contract_name.receive_name)
		99, 111, 110, 116, 114, 97, 99, 116, 95, 110, 97, 109, 101,
		46, 114, 101, 99, 101, 105, 118, 101, 95, 110, 97, 109, 101,
		0, 15, // params len
		8, 16, 0, 32, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, // params
	}
)

func Test_updateContractBody_Serialize(t *testing.T) {
	got, err := testUpdateContractBody.Serialize()
	if err != nil {
		t.Errorf("Serialize() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, testUpdateContractBodyBytes) {
		t.Errorf("Serialize() got = %v, want %v", got, testUpdateContractBodyBytes)
	}
}
