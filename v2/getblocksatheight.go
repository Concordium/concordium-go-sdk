package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBlocksAtHeight get a list of live blocks at a given height.
func (c *Client) GetBlocksAtHeight(ctx context.Context, req *pb.BlocksAtHeightRequest) (_ *pb.BlocksAtHeightResponse, err error) {
	blockAtHeight, err := c.grpcClient.GetBlocksAtHeight(ctx, req)
	if err != nil {
		return &pb.BlocksAtHeightResponse{}, err
	}

	return blockAtHeight, nil
}
