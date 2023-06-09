package v2

import (
	"context"
)

// GetBlockTransactionEvents get a list of transaction events in a given block.
// The stream will end when all the transaction events for a given block have been returned.
func (c *Client) GetBlockTransactionEvents(ctx context.Context, req *BlockHashInput) (_ Queries_GetBlockTransactionEventsClient, err error) {
	stream, err := c.grpcClient.GetBlockTransactionEvents(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
