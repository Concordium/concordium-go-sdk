package v2

import (
	"context"
)

// GetBlockItemStatus get the status of and information about a specific block item (transaction).
func (c *Client) GetBlockItemStatus(ctx context.Context, req *TransactionHash) (_ *BlockItemStatus, err error) {
	blockItemStatus, err := c.grpcClient.GetBlockItemStatus(ctx, req)
	if err != nil {
		return &BlockItemStatus{}, Error.Wrap(err)
	}

	return blockItemStatus, nil
}
