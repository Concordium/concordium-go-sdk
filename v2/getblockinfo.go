package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBlockInfo get information, such as height, timings, and transaction counts for the given block.
func (c *Client) GetBlockInfo(ctx context.Context, req *pb.BlockHashInput) (_ *pb.BlockInfo, err error) {
	blockInfo, err := c.grpcClient.GetBlockInfo(ctx, req)
	if err != nil {
		return &pb.BlockInfo{}, Error.Wrap(err)
	}

	return blockInfo, nil
}
