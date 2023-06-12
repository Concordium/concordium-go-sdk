package v2

import (
	"context"
)

// GetBlockPendingUpdates get the pending updates to chain parameters at the end of a given block.
// The stream will end when all the pending updates for a given block have been returned.
func (c *Client) GetBlockPendingUpdates(ctx context.Context, req *BlockHashInput) (_ Queries_GetBlockPendingUpdatesClient, err error) {
	stream, err := c.grpcClient.GetBlockPendingUpdates(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
