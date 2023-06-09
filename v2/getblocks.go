package v2

import (
	"context"
)

// GetBlocks returns a stream of blocks that arrive from the time
// the query is made onward. This can be used to listen for incoming blocks.
func (c *Client) GetBlocks(ctx context.Context) (_ Queries_GetBlocksClient, err error) {
	stream, err := c.grpcClient.GetBlocks(ctx, new(Empty))
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
