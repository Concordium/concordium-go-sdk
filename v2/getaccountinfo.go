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
		token := AccountToken{
			TokenID: TokenId{Value: t.GetTokenId().GetValue()},
			State: TokenAccountState{
				Balance: TokenAmount{Value: t.GetTokenAccountState().GetBalance().GetValue()},
			},
		}

		if t.GetTokenAccountState().GetModuleState() != nil {
			token.State.ModuleState = &RawCBOR{Bytes: t.GetTokenAccountState().GetModuleState().GetValue()}
		}

		info.Tokens = append(info.Tokens, token)
	}

	return info, nil
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
