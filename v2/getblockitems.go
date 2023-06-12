package v2

import (
	"context"
)

// GetBlockItems get the items of a block.
func (c *Client) GetBlockItems(ctx context.Context, req *BlockHashInput) (_ Queries_GetBlockItemsClient, err error) {
	stream, err := c.grpcClient.GetBlockItems(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
