package v2

import (
	"concordium-go-sdk/v2/pb"
	"context"
)

// GetBlocksAtHeight get a list of live blocks at a given height.
func (c *Client) GetBlocksAtHeight(ctx context.Context, req *pb.BlocksAtHeightRequest) (_ *pb.BlocksAtHeightResponse, err error) {
	blockAtHeight, err := c.grpcClient.GetBlocksAtHeight(ctx, req)
	if err != nil {
		return &pb.BlocksAtHeightResponse{}, Error.Wrap(err)
	}

	return blockAtHeight, nil
}
