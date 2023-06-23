package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAccountList retrieve the list of accounts that exist at the end of the given block.
func (c *Client) GetAccountList(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetAccountListClient, err error) {
	stream, err := c.grpcClient.GetAccountList(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
