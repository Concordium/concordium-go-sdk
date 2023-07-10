package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAccountInfo retrieve the information about the given account in the given block.
func (c *Client) GetAccountInfo(ctx context.Context, accId *pb.AccountIdentifierInput, b isBlockHashInput) (_ *pb.AccountInfo, err error) {
	accountInfo, err := c.GrpcClient.GetAccountInfo(ctx, &pb.AccountInfoRequest{
		BlockHash:         convertBlockHashInput(b),
		AccountIdentifier: accId,
	})
	if err != nil {
		return &pb.AccountInfo{}, err
	}

	return accountInfo, nil
}
