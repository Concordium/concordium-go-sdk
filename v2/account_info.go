package v2

import (
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

// AccountInfo represents the decoded account info returned by the node.
type AccountInfo struct {
	AccountNonce           Nonce
	AccountAmount          Amount
	AccountThreshold       AccountThreshold
	AccountEncryptedAmount AccountEncryptedAmount
	AccountIndex           AccountIndex
	AccountAddress         AccountAddress
	AvailableBalance       Amount
	Tokens                 []AccountToken
}

// AccountToken represents a token held by the account (PLT).
type AccountToken struct {
	TokenID TokenId
	State   TokenAccountState
}

// TokenModuleAccountState mirrors Rust SDK type.
type TokenModuleAccountState struct {
	AllowList  *bool                  `cbor:"allow_list,omitempty"`
	DenyList   *bool                  `cbor:"deny_list,omitempty"`
	Additional map[string]interface{} `cbor:",omitempty"`
}

// TokenAccountState represents account-level token info.
type TokenAccountState struct {
	Balance     TokenAmount
	ModuleState *RawCBOR // optional
}

func (s *TokenAccountState) DecodeModuleState() (*TokenModuleAccountState, error) {
	if s.ModuleState == nil || len(s.ModuleState.Bytes) == 0 {
		return &TokenModuleAccountState{}, nil
	}

	var decoded TokenModuleAccountState
	if err := cbor.Unmarshal(s.ModuleState.Bytes, &decoded); err != nil {
		return nil, fmt.Errorf("failed to decode module account state: %w", err)
	}

	return &decoded, nil
}
