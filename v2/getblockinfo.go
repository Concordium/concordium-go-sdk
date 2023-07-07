package v2

import (
	"context"
)

// GetBlockInfo get information, such as height, timings, and transaction counts for the given block.
func (c *Client) GetBlockInfo(ctx context.Context, req isBlockHashInput) (_ BlockInfo, err error) {
	blockInfo, err := c.GrpcClient.GetBlockInfo(ctx, convertBlockHashInput(req))
	if err != nil {
		return BlockInfo{}, err
	}

	return convertBlockInfo(blockInfo), nil
}
