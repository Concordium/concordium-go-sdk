package v2

import (
	"context"
)

// GetBlockSpecialEvents get a list of transaction events in a given block.
// The stream will end when all the transaction events for a given block have been returned.
func (c *Client) GetBlockSpecialEvents(ctx context.Context, req *BlockHashInput) (_ Queries_GetBlockSpecialEventsClient, err error) {
	stream, err := c.grpcClient.GetBlockSpecialEvents(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
