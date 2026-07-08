package v2

import (
	"context"
	"fmt"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetAccountInfo retrieves information about a given account in the given block.
func (c *Client) GetAccountInfo(ctx context.Context, accId *pb.AccountIdentifierInput, b isBlockHashInput) (*AccountInfo, error) {
	resp, err := c.GrpcClient.GetAccountInfo(ctx, &pb.AccountInfoRequest{
		BlockHash:         convertBlockHashInput(b),
		AccountIdentifier: accId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("empty response from GetAccountInfo")
	}

	info := &AccountInfo{
		AccountNonce:           Nonce{Value: resp.GetSequenceNumber().GetValue()},
		AccountAmount:          Amount{Value: resp.GetAmount().GetValue()},
		AccountAddress:         AccountAddress{Value: toFixed32(resp.GetAddress().GetValue())},
		AccountThreshold:       AccountThreshold{Value: uint8(resp.GetThreshold().GetValue())},
		AccountIndex:           AccountIndex{Value: resp.GetIndex().GetValue()},
		AvailableBalance:       Amount{Value: resp.GetAvailableBalance().GetValue()},
		AccountEncryptedAmount: convertEncryptedBalance(resp.GetEncryptedBalance()),
		Tokens:                 make([]AccountToken, 0, len(resp.GetTokens())),
	}

	for _, t := range resp.GetTokens() {
		state := t.GetTokenAccountState().GetBalance()
		token := AccountToken{
			TokenID: TokenId{Value: t.GetTokenId().GetValue()},
			State: TokenAccountState{
				Balance: TokenAmount{
					Value:    state.GetValue(),
					Decimals: uint8(state.GetDecimals()),
				},
			},
		}

		if t.GetTokenAccountState().GetModuleState() != nil {
			token.State.ModuleState = &RawCBOR{Bytes: t.GetTokenAccountState().GetModuleState().GetValue()}
		}

		info.Tokens = append(info.Tokens, token)
	}

	if creds := resp.GetCreds(); len(creds) > 0 {
		info.Credentials = make(map[CredentialIndex]CredentialInfo, len(creds))
		for idx, cred := range creds {
			info.Credentials[CredentialIndex(idx)] = convertCredential(cred)
		}
	}

	return info, nil
}

// convertCredential maps a pb.AccountCredential (initial or normal) into the
// SDK's CredentialInfo. All account verify keys are Ed25519.
func convertCredential(cred *pb.AccountCredential) CredentialInfo {
	var (
		credID *pb.CredentialRegistrationId
		pbKeys *pb.CredentialPublicKeys
	)
	switch {
	case cred.GetInitial() != nil:
		credID = cred.GetInitial().GetCredId()
		pbKeys = cred.GetInitial().GetKeys()
	case cred.GetNormal() != nil:
		credID = cred.GetNormal().GetCredId()
		pbKeys = cred.GetNormal().GetKeys()
	}

	info := CredentialInfo{
		CredID: CredentialRegistrationID{Value: credID.GetValue()},
		Keys: CredentialPublicKeys{
			Keys:      make(map[KeyIndex]VerifyKey, len(pbKeys.GetKeys())),
			Threshold: SignatureThreshold{Value: uint8(pbKeys.GetThreshold().GetValue())},
		},
	}
	for idx, vk := range pbKeys.GetKeys() {
		info.Keys.Keys[KeyIndex(idx)] = VerifyKey{Scheme: SchemeEd25519, Key: vk.GetEd25519Key()}
	}

	return info
}

// convertEncryptedBalance maps pb.EncryptedBalance → AccountEncryptedAmount.
func convertEncryptedBalance(pbBal *pb.EncryptedBalance) AccountEncryptedAmount {
	if pbBal == nil {
		return AccountEncryptedAmount{}
	}

	return AccountEncryptedAmount{
		SelfAmount: convertEncryptedAmount(pbBal.GetSelfAmount()),
	}
}

// convertEncryptedAmount maps pb.EncryptedAmount → EncryptedAmount.
func convertEncryptedAmount(pbAmt *pb.EncryptedAmount) *EncryptedAmount {
	if pbAmt == nil {
		return nil
	}
	return &EncryptedAmount{
		Value: pbAmt.GetValue(),
	}
}

// toFixed32 safely converts a byte slice to a [32]byte array.
func toFixed32(b []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], b)
	return arr
}
