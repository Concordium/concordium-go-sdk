package v2

import (
	"context"
)

// GetBlockInfo get information, such as height, timings, and transaction counts for the given block.
func (c *Client) GetBlockInfo(ctx context.Context, req *BlockHashInput) (_ *BlockInfo, err error) {
	blockInfo, err := c.grpcClient.GetBlockInfo(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return blockInfo, nil
}
