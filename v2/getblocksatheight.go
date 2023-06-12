package v2

import (
	"context"
)

// GetBlocksAtHeight get a list of live blocks at a given height.
func (c *Client) GetBlocksAtHeight(ctx context.Context, req *BlocksAtHeightRequest) (_ *BlocksAtHeightResponse, err error) {
	blockAtHeight, err := c.grpcClient.GetBlocksAtHeight(ctx, req)
	if err != nil {
		return &BlocksAtHeightResponse{}, Error.Wrap(err)
	}

	return blockAtHeight, nil
}
