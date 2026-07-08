package tests_test

import (
	"testing"

	v2 "github.com/Concordium/concordium-go-sdk/v2"

	"github.com/fxamacker/cbor/v2"
)

func TestTokenOperations_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		op        v2.TokenOperation
		expectKey string // expected top-level CBOR map key
	}{
		{
			name: "TransferOperation",
			op: &v2.TransferOperation{
				Transfer: v2.TokenTransfer{
					Amount:    v2.TokenAmount{Value: 42},
					Recipient: v2.CborHolderAccount{},
					Memo:      v2.CborMemo{},
				},
			},
			expectKey: "transfer",
		},
		{
			name:      "MintOperation",
			op:        &v2.MintOperation{Details: v2.TokenSupplyUpdateDetails{Amount: v2.TokenAmount{Value: 99}}},
			expectKey: "mint",
		},
		{
			name:      "BurnOperation",
			op:        &v2.BurnOperation{Details: v2.TokenSupplyUpdateDetails{Amount: v2.TokenAmount{Value: 77}}},
			expectKey: "burn",
		},
		{
			name:      "AddAllowListOperation",
			op:        &v2.AddAllowListOperation{Details: v2.TokenListUpdateDetails{}},
			expectKey: "add_allow_list",
		},
		{
			name:      "RemoveAllowListOperation",
			op:        &v2.RemoveAllowListOperation{Details: v2.TokenListUpdateDetails{}},
			expectKey: "remove_allow_list",
		},
		{
			name:      "AddDenyListOperation",
			op:        &v2.AddDenyListOperation{Details: v2.TokenListUpdateDetails{}},
			expectKey: "add_deny_list",
		},
		{
			name:      "RemoveDenyListOperation",
			op:        &v2.RemoveDenyListOperation{Details: v2.TokenListUpdateDetails{}},
			expectKey: "remove_deny_list",
		},
		{
			name:      "PauseOperation",
			op:        &v2.PauseOperation{Details: v2.TokenPauseDetails{}},
			expectKey: "pause",
		},
		{
			name:      "UnpauseOperation",
			op:        &v2.UnpauseOperation{Details: v2.TokenPauseDetails{}},
			expectKey: "unpause",
		},
		{
			name:      "UnknownOperation",
			op:        &v2.UnknownOperation{},
			expectKey: "", // empty map expected
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.op.EncodeCBOR()
			if err != nil {
				t.Fatalf("EncodeCBOR failed: %v", err)
			}

			var generic any
			if err := cbor.Unmarshal(data, &generic); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// fxamacker/cbor by default decodes into map[interface{}]interface{} for generic values.
			// Accept both map[string]any and map[interface{}]interface{}.
			if mStr, ok := generic.(map[string]any); ok {
				if tc.expectKey == "" {
					if len(mStr) != 0 {
						t.Fatalf("expected empty map for UnknownOperation, got: %#v", mStr)
					}
					return
				}
				if _, ok := mStr[tc.expectKey]; !ok {
					t.Fatalf("expected top-level key %q in CBOR map, got keys: %#v", tc.expectKey, mStr)
				}
				return
			}

			mAny, ok := generic.(map[interface{}]interface{})
			if !ok {
				t.Fatalf("expected map at top-level, got: %T", generic)
			}
			if tc.expectKey == "" {
				if len(mAny) != 0 {
					t.Fatalf("expected empty map for UnknownOperation, got: %#v", mAny)
				}
				return
			}

			found := false
			for k := range mAny {
				if ks, ok := k.(string); ok && ks == tc.expectKey {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected top-level key %q in CBOR map, got keys: %#v", tc.expectKey, mAny)
			}
		})
	}
}
