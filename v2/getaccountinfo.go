package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAccountInfo retrieve the information about the given account in the given block.
func (c *Client) GetAccountInfo(ctx context.Context, req *pb.AccountInfoRequest) (_ *pb.AccountInfo, err error) {
	accountInfo, err := c.grpcClient.GetAccountInfo(ctx, req)
	if err != nil {
		return &pb.AccountInfo{}, err
	}

	return accountInfo, nil
}
