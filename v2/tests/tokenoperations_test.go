package tests_test

import (
	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"reflect"
	"testing"

	"github.com/fxamacker/cbor/v2"
)

func TestTokenOperations_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		op   v2.TokenOperation
		new  func() v2.TokenOperation
	}{
		{
			name: "TransferOperation",
			op: &v2.TransferOperation{
				Transfer: v2.TokenTransfer{
					Amount:    v2.TokenAmount{Value: 42},
					Recipient: v2.CborTokenHolder{},
					Memo:      v2.CborMemo{},
				},
			},
			new: func() v2.TokenOperation { return &v2.TransferOperation{} },
		},
		{
			name: "MintOperation",
			op:   &v2.MintOperation{Details: v2.TokenSupplyUpdateDetails{Amount: v2.TokenAmount{Value: 99}}},
			new:  func() v2.TokenOperation { return &v2.MintOperation{} },
		},
		{
			name: "BurnOperation",
			op:   &v2.BurnOperation{Details: v2.TokenSupplyUpdateDetails{Amount: v2.TokenAmount{Value: 77}}},
			new:  func() v2.TokenOperation { return &v2.BurnOperation{} },
		},
		{
			name: "AddAllowListOperation",
			op:   &v2.AddAllowListOperation{Details: v2.TokenListUpdateDetails{}},
			new:  func() v2.TokenOperation { return &v2.AddAllowListOperation{} },
		},
		{
			name: "RemoveAllowListOperation",
			op:   &v2.RemoveAllowListOperation{Details: v2.TokenListUpdateDetails{}},
			new:  func() v2.TokenOperation { return &v2.RemoveAllowListOperation{} },
		},
		{
			name: "AddDenyListOperation",
			op:   &v2.AddDenyListOperation{Details: v2.TokenListUpdateDetails{}},
			new:  func() v2.TokenOperation { return &v2.AddDenyListOperation{} },
		},
		{
			name: "RemoveDenyListOperation",
			op:   &v2.RemoveDenyListOperation{Details: v2.TokenListUpdateDetails{}},
			new:  func() v2.TokenOperation { return &v2.RemoveDenyListOperation{} },
		},
		{
			name: "PauseOperation",
			op:   &v2.PauseOperation{Details: v2.TokenPauseDetails{}},
			new:  func() v2.TokenOperation { return &v2.PauseOperation{} },
		},
		{
			name: "UnpauseOperation",
			op:   &v2.UnpauseOperation{Details: v2.TokenPauseDetails{}},
			new:  func() v2.TokenOperation { return &v2.UnpauseOperation{} },
		},
		{
			name: "UnknownOperation",
			op:   &v2.UnknownOperation{},
			new:  func() v2.TokenOperation { return &v2.UnknownOperation{} },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.op.EncodeCBOR()
			if err != nil {
				t.Fatalf("EncodeCBOR failed: %v", err)
			}

			out := tc.new()
			if err := cbor.Unmarshal(data, out); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if !reflect.DeepEqual(tc.op, out) {
				t.Errorf("objects not equal after marshal/unmarshal\noriginal: %#v\nunmarshalled: %#v", tc.op, out)
			}
		})
	}
}
